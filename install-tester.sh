#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")" && pwd)"
env_file="$root/deploy/.env"
rendered="$root/deploy/nginx/rendered"
ai=""
calls=""
model=""
ai_memory=""
cloudflare=""
server=""
domain=""
proxy=""

for arg in "$@"; do
  case "$arg" in
    --ai) ai=1 ;;
    --no-ai) ai=0 ;;
    --calls) calls=1 ;;
    --no-calls) calls=0 ;;
    --model=*) model="${arg#--model=}" ;;
    --server) server=1 ;;
    --domain=*) server=1; domain="${arg#--domain=}" ;;
    --proxy=docker|--proxy=host) server=1; proxy="${arg#--proxy=}" ;;
    --cloudflare) server=1; cloudflare=1 ;;
    --no-cloudflare) cloudflare=0 ;;
    --ai-memory=*) ai_memory="${arg#--ai-memory=}" ;;
    --local) server=0 ;;
    -h|--help)
      echo "Использование: ./install-tester.sh [--ai|--no-ai] [--calls|--no-calls] [--model=qwen2.5:3b] [--ai-memory=5g] [--server] [--domain=lodestar.example.kz] [--proxy=docker|host] [--cloudflare] [--local]"
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

udp_busy() {
  command -v ss >/dev/null 2>&1 && ss -lunH "( sport >= :$1 and sport <= :$2 )" 2>/dev/null | grep -q .
}

port_ours() {
  docker ps -q --filter "label=com.docker.compose.project=$project" --filter "publish=$1" 2>/dev/null | grep -q .
}

port_taken() {
  port_busy "$1" && ! port_ours "$1"
}

pick_port() {
  local key="$1" port
  port="$(get_key "$key")"
  [ -n "$port" ] || port="$2"
  while port_taken "$port"; do
    port=$((port + 1))
  done
  if [ "$port" != "$(get_key "$key")" ] && [ "$port" != "$2" ]; then
    warn "Порт $2 занят, $key = $port."
  fi
  set_key "$key" "$port"
  picked="$port"
}

check_resources() {
  local mem_kb disk_kb
  if [ -r /proc/meminfo ]; then
    mem_kb="$(awk '/MemAvailable/ {print $2}' /proc/meminfo)"
    if [ "$ai" = 1 ] && [ "$mem_kb" -lt 7000000 ]; then
      warn "Свободно памяти $((mem_kb / 1024 / 1024)) ГБ. С ИИ лучше от 7 ГБ свободных: ИИ берёт до $ai_memory, остальное сервер, база и звонки. Если сервер начнёт тормозить, запустите с --model=qwen2.5:1.5b или --no-ai."
    elif [ "$mem_kb" -lt 2000000 ]; then
      warn "Свободно памяти $((mem_kb / 1024 / 1024)) ГБ. Без ИИ нужно от 2 ГБ свободных."
    fi
  fi
  disk_kb="$(df -Pk "$root" | awk 'NR == 2 {print $4}')"
  if [ "$disk_kb" -lt 20000000 ]; then
    warn "Свободно на диске $((disk_kb / 1024 / 1024)) ГБ. Нужно около 20 ГБ: образы, веб-сборка и модели ИИ."
  fi
  if [ -z "$(get_key LODESTAR_LLM_URL)" ] || [ "$(get_key LODESTAR_LLM_URL)" = "http://ollama:11434/v1" ]; then
    case "$model" in
      qwen2.5:0.5b|qwen2.5:1.5b|qwen2.5:3b|llama3.2:1b|llama3.2:3b|gemma2:2b) ;;
      *) warn "Модель $model может не уложиться в $ai_memory памяти, отведённые ИИ. Проверенная модель: qwen2.5:3b." ;;
    esac
  fi
}

real_ip() {
  local ips
  [ "$cloudflare" = 1 ] || return 0
  ips="$( { curl -fsS --max-time 5 https://www.cloudflare.com/ips-v4 && echo && curl -fsS --max-time 5 https://www.cloudflare.com/ips-v6; } 2>/dev/null || true)"
  if ! printf '%s\n' "$ips" | grep -q '/'; then
    ips="$(cat "$root/deploy/nginx/cloudflare-ips.txt")"
  fi
  printf '%s\n' "$ips" | awk 'NF {printf "    set_real_ip_from %s;\n", $1}'
  echo "    real_ip_header CF-Connecting-IP;"
}

render_nginx() {
  local upstream_server upstream_s3 upstream_rtc upstream_mail certs htpasswd
  if [ "$proxy" = docker ]; then
    upstream_server="server:8080"; upstream_s3="minio:9000"; upstream_rtc="livekit:7880"; upstream_mail="mailpit:8025"
    htpasswd="/etc/nginx/htpasswd"
  else
    upstream_server="127.0.0.1:$http_port"; upstream_s3="127.0.0.1:$s3_port"; upstream_rtc="127.0.0.1:$rtc_port"; upstream_mail="127.0.0.1:$mail_port"
    htpasswd="$rendered/htpasswd"
  fi
  certs="/etc/letsencrypt/live/$domain"
  mkdir -p "$rendered"
  local templates=("$root/deploy/nginx/lodestar.conf.template")
  [ "$calls" = 1 ] && templates+=("$root/deploy/nginx/rtc.conf.template")
  real_ip > "$rendered/realip.part"
  cat "${templates[@]}" | sed -e "s|__NAMES__|${names[*]}|g" -e "s|__DOMAIN__|$domain|g" -e "s|__S3_HOST__|$s3_host|g" \
    -e "s|__MAIL_HOST__|$mail_host|g" -e "s|__RTC_HOST__|$rtc_host|g" -e "s|__SERVER__|$upstream_server|g" \
    -e "s|__S3__|$upstream_s3|g" -e "s|__RTC__|$upstream_rtc|g" -e "s|__MAIL__|$upstream_mail|g" -e "s|__CERTS__|$certs|g" \
    -e "s|__HTPASSWD__|$htpasswd|g" -e "/__REALIP__/{r $rendered/realip.part" -e "d}" > "$rendered/lodestar.conf"
  rm -f "$rendered/realip.part"
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
    --register-unsafely-without-email $(printf -- '-d %s ' "${names[@]}"); then
    warn "Сертификат не выдан. Проверьте, что порт 80 открыт снаружи и адрес $domain указывает на этот сервер, и запустите установщик снова."
    exit 1
  fi
}

say "Настройки в deploy/.env"
touch "$env_file"
if [ -z "$server" ]; then
  if has_key LODESTAR_DOMAIN; then server=1; else server=0; fi
fi
[ -n "$ai" ] || ai="$(get_key LODESTAR_AI)"
[ -n "$ai" ] || ai=1
set_key LODESTAR_AI "$ai"
[ -n "$model" ] || model="$(get_key LODESTAR_LLM_MODEL)"
[ -n "$model" ] || model="qwen2.5:3b"
[ -n "$ai_memory" ] || ai_memory="$(get_key LODESTAR_AI_MEMORY)"
if [ -z "$ai_memory" ]; then
  case "$model" in
    qwen2.5:0.5b|qwen2.5:1.5b|llama3.2:1b) ai_memory=3g ;;
    *) ai_memory=5g ;;
  esac
fi
set_key LODESTAR_AI_MEMORY "$ai_memory"
keep_key LODESTAR_SECRET_KEY "$(secret 32)"
keep_key POSTGRES_PASSWORD "$(word 24)"
keep_key MINIO_ROOT_USER lodestar
keep_key MINIO_ROOT_PASSWORD "$(word 24)"
keep_key LIVEKIT_API_KEY lodestar
keep_key LIVEKIT_API_SECRET "$(word 40)"

if [ "$server" = 1 ]; then
  keep_key COMPOSE_PROJECT_NAME lodestar
elif docker volume ls -q 2>/dev/null | grep -q '^deploy_'; then
  keep_key COMPOSE_PROJECT_NAME deploy
else
  keep_key COMPOSE_PROJECT_NAME lodestar
fi
project="$(get_key COMPOSE_PROJECT_NAME)"
files=(-p "$project" -f "$root/deploy/compose.yaml" -f "$root/deploy/tester.yaml")

if [ "$server" = 1 ]; then
  [ -n "$calls" ] || calls="$(get_key LODESTAR_CALLS)"
  [ -n "$calls" ] || calls=0
  set_key LODESTAR_CALLS "$calls"
else
  drop_key LODESTAR_CALLS
  [ -n "$calls" ] || calls=1
fi

pick_port LODESTAR_HTTP_PORT 8080; http_port="$picked"
pick_port LODESTAR_S3_PORT 9000; s3_port="$picked"
pick_port MAILPIT_PORT 8025; mail_port="$picked"
pick_port LIVEKIT_PORT 7880; rtc_port="$picked"
if [ "$calls" = 1 ]; then
  if port_taken 7881 || { udp_busy 50000 50200 && ! port_ours 7881; }; then
    warn "Порт 7881/tcp или диапазон 50000–50200/udp уже занят. Звонки могут не заработать; без них: --no-calls."
  fi
fi

if [ "$server" = 1 ]; then
  check_resources
  ip="$(public_ip)"
  [ -n "$domain" ] || domain="$(get_key LODESTAR_DOMAIN)"
  [ -n "$domain" ] || domain="${ip//./-}.sslip.io"
  [ -n "$cloudflare" ] || cloudflare="$(get_key LODESTAR_CLOUDFLARE)"
  [ -n "$cloudflare" ] || cloudflare=0
  set_key LODESTAR_CLOUDFLARE "$cloudflare"
  if [ "$cloudflare" = 1 ]; then
    s3_host="${domain%%.*}-s3.${domain#*.}"; mail_host="${domain%%.*}-mail.${domain#*.}"; rtc_host="${domain%%.*}-rtc.${domain#*.}"
  else
    s3_host="s3.$domain"; mail_host="mail.$domain"; rtc_host="rtc.$domain"
  fi
  names=("$domain" "$s3_host" "$mail_host")
  [ "$calls" = 1 ] && names+=("$rtc_host")
  [ -n "$proxy" ] || proxy="$(get_key LODESTAR_PROXY)"
  if [ -z "$proxy" ]; then
    proxy=docker
    for port in 80 443; do
      if port_taken "$port"; then
        proxy=host
      fi
    done
    if [ "$proxy" = host ]; then
      warn "Порты 80 и 443 уже заняты, поэтому Lodestar не запускает свой nginx и подключается к вашему (режим --proxy=host)."
    fi
  fi
  set_key LODESTAR_DOMAIN "$domain"
  set_key LODESTAR_PROXY "$proxy"
  set_key LODESTAR_HOST_IP "$ip"
  set_key LODESTAR_PUBLIC_URL "https://$domain"
  set_key LODESTAR_S3_PUBLIC_ENDPOINT "https://$s3_host"
  keep_key MAILPIT_PASSWORD "$(word 16)"
  mail_password="$(get_key MAILPIT_PASSWORD)"
  files+=(-f "$root/deploy/vds.yaml")
  render_nginx
  if [ "$proxy" = docker ]; then
    for port in 80 443; do
      if port_taken "$port"; then
        warn "Порт $port уже занят другой программой. Освободите его или запустите с --proxy=host, чтобы подключить Lodestar к своему nginx."
      fi
    done
  fi
else
  ip="$(host_ip)"
  drop_key LODESTAR_DOMAIN
  drop_key LODESTAR_PROXY
  set_key LODESTAR_HOST_IP "$ip"
  set_key LODESTAR_PUBLIC_URL "http://$ip:$http_port"
  set_key LODESTAR_S3_PUBLIC_ENDPOINT "http://$ip:$s3_port"
fi

profiles=""
if [ "$ai" = 1 ]; then
  profiles="ai"
  llm_url="$(get_key LODESTAR_LLM_URL)"
  [ -n "$llm_url" ] || llm_url="http://ollama:11434/v1"
  set_key LODESTAR_LLM_URL "$llm_url"
  set_key LODESTAR_LLM_MODEL "$model"
  if [ "$llm_url" = "http://ollama:11434/v1" ]; then
    set_key LODESTAR_OLLAMA_LLM_MODEL "$model"
  else
    set_key LODESTAR_OLLAMA_LLM_MODEL ""
  fi
  set_key LODESTAR_EMBED_URL "http://ollama:11434/v1"
  set_key LODESTAR_EMBED_MODEL "bge-m3"
else
  set_key LODESTAR_LLM_URL ""
  set_key LODESTAR_EMBED_URL ""
fi
if [ "$calls" = 1 ]; then
  profiles="${profiles:+$profiles,}calls"
  if [ "$server" = 1 ]; then
    set_key LODESTAR_LIVEKIT_URL "wss://$rtc_host"
  else
    set_key LODESTAR_LIVEKIT_URL "ws://$ip:$rtc_port"
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
  if curl -fsS "http://localhost:$http_port/healthz" >/dev/null 2>&1; then
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
  Письма:          https://$mail_host  (логин tester, пароль $mail_password)
INFO
if [ "$proxy" = host ]; then
  cat <<STEPS

Подключение к nginx сервера (скрипт этого не делает):
  1. DNS: A-записи ${names[*]} на $ip
  2. sudo certbot certonly --nginx --cert-name $domain $(printf -- '-d %s ' "${names[@]}")
  3. sudo ln -s $rendered/lodestar.conf /etc/nginx/conf.d/lodestar.conf
  4. sudo nginx -t && sudo systemctl reload nginx
STEPS
  case "$root" in
    /root/*|/home/*) warn "Проект лежит в $root. nginx может не прочитать оттуда пароль почты ($rendered/htpasswd). Надёжнее держать его, например, в /opt/lodestar." ;;
  esac
fi
if [ "$cloudflare" = 1 ]; then
  warn "Cloudflare: имена плоские (${names[*]}), их покрывает бесплатный сертификат Cloudflare *.${domain#*.}. Реальные IP берутся из CF-Connecting-IP только от адресов Cloudflare. Сертификат выпускайте при выключенном проксировании (серое облако), потом включайте и ставьте SSL Full (strict). «Always Use HTTPS» для этих имён не включайте, иначе не пройдёт продление: на https nginx переводит сам. Бесплатный тариф пропускает загрузки до 100 МБ."
fi
if [ "$calls" = 1 ]; then
  warn "Настройки сервера скрипт не меняет. Снаружи должны быть открыты порты 80 и 443, для звонков ещё 7881/tcp и 50000–50200/udp. Если включён файрвол (ufw, firewalld или у провайдера), откройте их там."
else
  warn "Настройки сервера скрипт не меняет. Снаружи должны быть открыты порты 80 и 443. Звонки выключены, включить: --calls (нужны ещё 7881/tcp и 50000–50200/udp)."
fi
else
cat <<INFO

Готово.
  Lodestar:        http://localhost:$http_port  (с телефона в той же Wi-Fi: http://$ip:$http_port)
  Код вуза:        demo
  Вход:            раздел «Режим тестировщика» на экране входа
  Письма:          http://localhost:$mail_port
INFO
fi
if [ "$ai" = 1 ]; then
  echo "  ИИ-навигатор:    модель $model ($llm_url), эмбеддинги bge-m3; загрузка локальных моделей: docker compose ${files[*]} --env-file $env_file logs -f ollama-models"
fi
echo "  Остановить:      docker compose ${files[*]} --env-file $env_file down  (./install-tester.sh снова обновит и запустит)"
if [ "${#warnings[@]}" -gt 0 ]; then
  printf '\nПредупреждения (%d):\n' "${#warnings[@]}"
  printf '  - %s\n' "${warnings[@]}"
fi
