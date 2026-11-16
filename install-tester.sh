#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")" && pwd)"
env_file="$root/deploy/.env"
rendered="$root/deploy/nginx/rendered"
ai=1
calls=1
model="qwen2.5:3b"
server=""
domain=""
proxy=""

for arg in "$@"; do
  case "$arg" in
    --no-ai) ai=0 ;;
    --no-calls) calls=0 ;;
    --model=*) model="${arg#--model=}" ;;
    --server) server=1 ;;
    --domain=*) server=1; domain="${arg#--domain=}" ;;
    --proxy=docker|--proxy=host) server=1; proxy="${arg#--proxy=}" ;;
    --local) server=0 ;;
    -h|--help)
      echo "Использование: ./install-tester.sh [--no-ai] [--no-calls] [--model=qwen2.5:3b] [--server] [--domain=lodestar.example.kz] [--proxy=docker|host] [--local]"
      exit 0 ;;
    *) echo "Неизвестный параметр: $arg" >&2; exit 1 ;;
  esac
done

say() { printf '\n==> %s\n' "$1"; }
warnings=()
warn() {
  warnings+=("$1")
  printf '\n[!] %s\n' "$1" >&2
}

if ! command -v docker >/dev/null 2>&1; then
  echo "Нужен Docker: https://docs.docker.com/get-docker/ . Поставьте его, запустите и повторите." >&2
  exit 1
fi
if ! docker info >/dev/null 2>&1; then
  echo "Docker установлен, но не запущен или нет прав на него. Запустите Docker и повторите." >&2
  exit 1
fi
if ! docker compose version >/dev/null 2>&1; then
  echo "Нужен Docker Compose v2 (команда docker compose)." >&2
  exit 1
fi

host_ip() {
  local ip=""
  if command -v ip >/dev/null 2>&1; then
    ip="$(ip route get 1.1.1.1 2>/dev/null | awk '{for (i = 1; i < NF; i++) if ($i == "src") print $(i + 1)}' | head -n1)"
  fi
  if [ -z "$ip" ] && command -v ipconfig >/dev/null 2>&1; then
    ip="$(ipconfig getifaddr en0 2>/dev/null || ipconfig getifaddr en1 2>/dev/null || true)"
  fi
  echo "${ip:-127.0.0.1}"
}

secret() {
  head -c "$1" /dev/urandom | base64 | tr -d '\n'
}

word() {
  head -c 64 /dev/urandom | base64 | tr -dc 'A-Za-z0-9' | head -c "$1"
}

has_key() {
  grep -q "^$1=." "$env_file"
}

get_key() {
  grep "^$1=" "$env_file" | cut -d= -f2- || true
}

set_key() {
  local tmp
  tmp="$(mktemp)"
  grep -v "^$1=" "$env_file" > "$tmp" || true
  printf '%s=%s\n' "$1" "$2" >> "$tmp"
  mv "$tmp" "$env_file"
}

drop_key() {
  local tmp
  tmp="$(mktemp)"
  grep -v "^$1=" "$env_file" > "$tmp" || true
  mv "$tmp" "$env_file"
}

keep_key() {
  has_key "$1" || set_key "$1" "$2"
}

public_ip() {
  curl -fsS4 --max-time 5 https://api.ipify.org 2>/dev/null || host_ip
}

port_busy() {
  command -v ss >/dev/null 2>&1 && ss -ltnH "( sport = :$1 )" 2>/dev/null | grep -q .
}

check_resources() {
  local mem_kb disk_kb
  if [ -r /proc/meminfo ]; then
    mem_kb="$(awk '/MemTotal/ {print $2}' /proc/meminfo)"
    if [ "$ai" = 1 ] && [ "$mem_kb" -lt 7500000 ]; then
      warn "Памяти $((mem_kb / 1024 / 1024)) ГБ. С ИИ нужно от 8 ГБ: ИИ берёт до 5 ГБ, остальное сервер, база и звонки. Можно запустить с --no-ai."
    elif [ "$mem_kb" -lt 3500000 ]; then
      warn "Памяти $((mem_kb / 1024 / 1024)) ГБ. Без ИИ нужно от 4 ГБ."
    fi
  fi
  disk_kb="$(df -Pk "$root" | awk 'NR == 2 {print $4}')"
  if [ "$disk_kb" -lt 20000000 ]; then
    warn "Свободно на диске $((disk_kb / 1024 / 1024)) ГБ. Нужно около 20 ГБ: образы, веб-сборка и модели ИИ."
  fi
  case "$model" in
    qwen2.5:0.5b|qwen2.5:1.5b|qwen2.5:3b|llama3.2:1b|llama3.2:3b|gemma2:2b) ;;
    *) warn "Модель $model может не уложиться в 5 ГБ памяти, отведённые ИИ. Проверенная модель: qwen2.5:3b." ;;
  esac
}

render_nginx() {
  local upstream_server upstream_s3 upstream_rtc upstream_mail certs htpasswd
  if [ "$proxy" = docker ]; then
    upstream_server="server:8080"; upstream_s3="minio:9000"; upstream_rtc="livekit:7880"; upstream_mail="mailpit:8025"
    htpasswd="/etc/nginx/htpasswd"
  else
    upstream_server="127.0.0.1:8080"; upstream_s3="127.0.0.1:9000"; upstream_rtc="127.0.0.1:7880"; upstream_mail="127.0.0.1:8025"
    htpasswd="$rendered/htpasswd"
  fi
  certs="/etc/letsencrypt/live/$domain"
  mkdir -p "$rendered"
  sed -e "s|__DOMAIN__|$domain|g" -e "s|__SERVER__|$upstream_server|g" -e "s|__S3__|$upstream_s3|g" \
    -e "s|__RTC__|$upstream_rtc|g" -e "s|__MAIL__|$upstream_mail|g" -e "s|__CERTS__|$certs|g" \
    -e "s|__HTPASSWD__|$htpasswd|g" "$root/deploy/nginx/lodestar.conf.template" > "$rendered/lodestar.conf"
  if [ ! -s "$rendered/htpasswd" ]; then
    docker run --rm httpd:2.4-alpine htpasswd -nbB tester "$mail_password" > "$rendered/htpasswd"
  fi
}

issue_certificate() {
  local compose=(docker compose "${files[@]}" --env-file "$env_file")
  if "${compose[@]}" run --rm --no-deps --entrypoint sh certbot -c "test -f /etc/letsencrypt/live/$domain/fullchain.pem" >/dev/null 2>&1; then
    return
  fi
  say "Получаю сертификат HTTPS для $domain"
  "${compose[@]}" stop nginx >/dev/null 2>&1 || true
  if ! "${compose[@]}" run --rm --no-deps -p 80:80 --entrypoint certbot certbot certonly --standalone -n --agree-tos \
    --register-unsafely-without-email -d "$domain" -d "s3.$domain" -d "rtc.$domain" -d "mail.$domain"; then
    warn "Сертификат не выдан. Проверьте, что порт 80 открыт снаружи и адрес $domain указывает на этот сервер, и запустите установщик снова."
    exit 1
  fi
}

say "Настройки в deploy/.env"
touch "$env_file"
if [ -z "$server" ]; then
  if has_key LODESTAR_DOMAIN; then server=1; else server=0; fi
fi
keep_key LODESTAR_SECRET_KEY "$(secret 32)"
keep_key POSTGRES_PASSWORD "$(word 24)"
keep_key MINIO_ROOT_USER lodestar
keep_key MINIO_ROOT_PASSWORD "$(word 24)"
keep_key LIVEKIT_API_KEY lodestar
keep_key LIVEKIT_API_SECRET "$(word 40)"

files=(-f "$root/deploy/compose.yaml" -f "$root/deploy/tester.yaml")
if [ "$server" = 1 ]; then
  check_resources
  ip="$(public_ip)"
  [ -n "$domain" ] || domain="$(get_key LODESTAR_DOMAIN)"
  [ -n "$domain" ] || domain="${ip//./-}.sslip.io"
  [ -n "$proxy" ] || proxy="$(get_key LODESTAR_PROXY)"
  [ -n "$proxy" ] || proxy=docker
  set_key LODESTAR_DOMAIN "$domain"
  set_key LODESTAR_PROXY "$proxy"
  set_key LODESTAR_HOST_IP "$ip"
  set_key LODESTAR_PUBLIC_URL "https://$domain"
  set_key LODESTAR_S3_PUBLIC_ENDPOINT "https://s3.$domain"
  keep_key MAILPIT_PASSWORD "$(word 16)"
  mail_password="$(get_key MAILPIT_PASSWORD)"
  files+=(-f "$root/deploy/vds.yaml")
  render_nginx
  if [ "$proxy" = docker ]; then
    for port in 80 443; do
      if port_busy "$port" && ! docker compose "${files[@]}" --env-file "$env_file" ps --status running nginx 2>/dev/null | grep -q nginx; then
        warn "Порт $port уже занят другой программой. Освободите его или запустите с --proxy=host, чтобы подключить Lodestar к своему nginx."
      fi
    done
  fi
else
  ip="$(host_ip)"
  drop_key LODESTAR_DOMAIN
  drop_key LODESTAR_PROXY
  set_key LODESTAR_HOST_IP "$ip"
  set_key LODESTAR_PUBLIC_URL "http://$ip:8080"
  set_key LODESTAR_S3_PUBLIC_ENDPOINT "http://$ip:9000"
fi

profiles=""
if [ "$ai" = 1 ]; then
  profiles="ai"
  set_key LODESTAR_LLM_URL "http://ollama:11434/v1"
  set_key LODESTAR_LLM_MODEL "$model"
  set_key LODESTAR_EMBED_URL "http://ollama:11434/v1"
  set_key LODESTAR_EMBED_MODEL "bge-m3"
else
  set_key LODESTAR_LLM_URL ""
  set_key LODESTAR_EMBED_URL ""
fi
if [ "$calls" = 1 ]; then
  profiles="${profiles:+$profiles,}calls"
  if [ "$server" = 1 ]; then
    set_key LODESTAR_LIVEKIT_URL "wss://rtc.$domain"
  else
    set_key LODESTAR_LIVEKIT_URL "ws://$ip:7880"
  fi
else
  set_key LODESTAR_LIVEKIT_URL ""
fi
if [ "$server" = 1 ] && [ "$proxy" = docker ]; then
  profiles="${profiles:+$profiles,}proxy"
fi
set_key COMPOSE_PROFILES "$profiles"

if [ "$server" = 1 ] && [ "$proxy" = docker ]; then
  issue_certificate
fi

say "Сборка и запуск. Первый раз 15–30 минут: собирается веб-версия и скачиваются модели"
docker compose "${files[@]}" --env-file "$env_file" up -d --build --remove-orphans

say "Жду сервер"
ready=0
for _ in $(seq 1 90); do
  if curl -fsS "http://localhost:8080/healthz" >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 2
done
[ "$ready" = 1 ] || warn "Сервер не ответил за 3 минуты. Журнал: docker compose ${files[*]} --env-file $env_file logs server"

if [ "$server" = 1 ]; then
cat <<INFO

Готово.
  Lodestar:        https://$domain  (веб-версия; в приложении: «Сменить вуз», адрес https://$domain)
  Код вуза:        demo
  Вход:            раздел «Режим тестировщика» на экране входа
  Письма:          https://mail.$domain  (логин tester, пароль $mail_password)
INFO
if [ "$proxy" = host ]; then
  warn "Lodestar не занимает порты 80 и 443. Подключите его к своему nginx: include $rendered/lodestar.conf, а сертификат на $domain, s3.$domain, rtc.$domain и mail.$domain выпустите сами (certbot --nginx или свой)."
fi
warn "Настройки сервера скрипт не меняет. Снаружи должны быть открыты порты 80 и 443, для звонков ещё 7881/tcp и 50000–50200/udp. Если включён файрвол (ufw, firewalld или у провайдера), откройте их там."
else
cat <<INFO

Готово.
  Lodestar:        http://localhost:8080  (с телефона в той же Wi-Fi: http://$ip:8080)
  Код вуза:        demo
  Вход:            раздел «Режим тестировщика» на экране входа
  Письма:          http://localhost:8025
INFO
fi
if [ "$ai" = 1 ]; then
  echo "  ИИ-навигатор:    модели $model и bge-m3 (вместе около 3 ГБ) качаются в фоне, следить: docker compose ${files[*]} --env-file $env_file logs -f ollama-models"
fi
echo "  Остановить:      docker compose ${files[*]} --env-file $env_file down  (./install-tester.sh снова обновит и запустит)"
if [ "${#warnings[@]}" -gt 0 ]; then
  printf '\nПредупреждения (%d):\n' "${#warnings[@]}"
  printf '  - %s\n' "${warnings[@]}"
fi
