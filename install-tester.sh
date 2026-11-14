#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")" && pwd)"
env_file="$root/deploy/.env"
ai=1
calls=1
model="qwen2.5:3b"

for arg in "$@"; do
  case "$arg" in
    --no-ai) ai=0 ;;
    --no-calls) calls=0 ;;
    --model=*) model="${arg#--model=}" ;;
    -h|--help)
      echo "Использование: ./install-tester.sh [--no-ai] [--no-calls] [--model=qwen2.5:7b]"
      exit 0 ;;
    *) echo "Неизвестный параметр: $arg" >&2; exit 1 ;;
  esac
done

say() { printf '\n==> %s\n' "$1"; }

if ! command -v docker >/dev/null 2>&1; then
  echo "Нужен Docker: https://docs.docker.com/get-docker/ . Поставьте его, запустите и повторите." >&2
  exit 1
fi
if ! docker info >/dev/null 2>&1; then
  echo "Docker установлен, но не запущен. Запустите Docker Desktop и повторите." >&2
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

set_key() {
  local tmp
  tmp="$(mktemp)"
  grep -v "^$1=" "$env_file" > "$tmp" || true
  printf '%s=%s\n' "$1" "$2" >> "$tmp"
  mv "$tmp" "$env_file"
}

keep_key() {
  has_key "$1" || set_key "$1" "$2"
}

say "Настройки в deploy/.env"
touch "$env_file"
keep_key LODESTAR_SECRET_KEY "$(secret 32)"
keep_key POSTGRES_PASSWORD "$(word 24)"
keep_key MINIO_ROOT_USER lodestar
keep_key MINIO_ROOT_PASSWORD "$(word 24)"
keep_key LIVEKIT_API_KEY lodestar
keep_key LIVEKIT_API_SECRET "$(word 40)"

ip="$(host_ip)"
set_key LODESTAR_HOST_IP "$ip"
set_key LODESTAR_PUBLIC_URL "http://$ip:8080"
set_key LODESTAR_S3_PUBLIC_ENDPOINT "http://$ip:9000"

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
  set_key LODESTAR_LIVEKIT_URL "ws://$ip:7880"
else
  set_key LODESTAR_LIVEKIT_URL ""
fi
set_key COMPOSE_PROFILES "$profiles"

say "Сборка и запуск. Первый раз 15–30 минут: собирается веб-версия и скачиваются модели"
docker compose -f "$root/deploy/compose.yaml" -f "$root/deploy/tester.yaml" --env-file "$env_file" up -d --build

say "Жду сервер"
for _ in $(seq 1 90); do
  if curl -fsS "http://localhost:8080/healthz" >/dev/null 2>&1; then
    break
  fi
  sleep 2
done

cat <<INFO

Готово.
  Lodestar:        http://localhost:8080  (с телефона в той же Wi-Fi: http://$ip:8080)
  Код вуза:        demo
  Вход:            раздел «Режим тестировщика» на экране входа
  Письма:          http://localhost:8025
INFO
if [ "$ai" = 1 ]; then
  echo "  ИИ-навигатор:    модели $model и bge-m3 качаются в фоне, следить: docker compose -f deploy/compose.yaml -f deploy/tester.yaml --env-file deploy/.env logs -f ollama-models"
fi
echo "  Остановить:      make tester-down  (или ./install-tester.sh снова, чтобы обновить и запустить)"
