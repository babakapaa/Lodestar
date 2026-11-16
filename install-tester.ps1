param(
  [switch]$NoAi,
  [switch]$NoCalls,
  [string]$Model = "qwen2.5:3b"
)

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$envFile = Join-Path $root "deploy\.env"

function Say($text) { Write-Host "`n==> $text" -ForegroundColor Cyan }

if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
  Write-Host "Нужен Docker Desktop: https://docs.docker.com/desktop/setup/install/windows-install/ . Поставьте его, запустите и повторите." -ForegroundColor Red
  exit 1
}
& cmd /c "docker info >nul 2>&1"
if ($LASTEXITCODE -ne 0) {
  Write-Host "Docker установлен, но не запущен. Запустите Docker Desktop и повторите." -ForegroundColor Red
  exit 1
}

function Get-HostIp {
  $route = Get-NetRoute -DestinationPrefix "0.0.0.0/0" -ErrorAction SilentlyContinue | Sort-Object RouteMetric | Select-Object -First 1
  if ($route) {
    $address = Get-NetIPAddress -InterfaceIndex $route.InterfaceIndex -AddressFamily IPv4 -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($address) { return $address.IPAddress }
  }
  return "127.0.0.1"
}

function New-Secret([int]$bytes) {
  $buffer = New-Object byte[] $bytes
  [System.Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($buffer)
  return [Convert]::ToBase64String($buffer)
}

function New-Word([int]$length) {
  return ((New-Secret 64) -replace "[^A-Za-z0-9]", "").Substring(0, $length)
}

$values = [ordered]@{}
if (Test-Path $envFile) {
  foreach ($line in Get-Content $envFile) {
    if ($line -match "^([A-Za-z0-9_]+)=(.*)$") { $values[$Matches[1]] = $Matches[2] }
  }
}

function Keep($key, $value) { if (-not $values[$key]) { $values[$key] = $value } }

if (-not $PSBoundParameters.ContainsKey('Model') -and $values['LODESTAR_LLM_MODEL']) {
  $Model = $values['LODESTAR_LLM_MODEL']
}

Say "Настройки в deploy\.env"
Keep "LODESTAR_SECRET_KEY" (New-Secret 32)
Keep "POSTGRES_PASSWORD" (New-Word 24)
Keep "MINIO_ROOT_USER" "lodestar"
Keep "MINIO_ROOT_PASSWORD" (New-Word 24)
Keep "LIVEKIT_API_KEY" "lodestar"
Keep "LIVEKIT_API_SECRET" (New-Word 40)

$ip = Get-HostIp
$values["LODESTAR_HOST_IP"] = $ip
$values["LODESTAR_PUBLIC_URL"] = "http://${ip}:8080"
$values["LODESTAR_S3_PUBLIC_ENDPOINT"] = "http://${ip}:9000"

$profiles = @()
if ($NoAi) {
  $values["LODESTAR_LLM_URL"] = ""
  $values["LODESTAR_EMBED_URL"] = ""
} else {
  $profiles += "ai"
  if (-not $values["LODESTAR_LLM_URL"]) { $values["LODESTAR_LLM_URL"] = "http://ollama:11434/v1" }
  $values["LODESTAR_LLM_MODEL"] = $Model
  $values["LODESTAR_OLLAMA_LLM_MODEL"] = if ($values["LODESTAR_LLM_URL"] -eq "http://ollama:11434/v1") { $Model } else { "" }
  $values["LODESTAR_EMBED_URL"] = "http://ollama:11434/v1"
  $values["LODESTAR_EMBED_MODEL"] = "bge-m3"
}
if ($NoCalls) {
  $values["LODESTAR_LIVEKIT_URL"] = ""
} else {
  $profiles += "calls"
  $values["LODESTAR_LIVEKIT_URL"] = "ws://${ip}:7880"
}
$values["COMPOSE_PROFILES"] = $profiles -join ","

$lines = foreach ($key in $values.Keys) { "$key=$($values[$key])" }
[System.IO.File]::WriteAllLines($envFile, [string[]]$lines, (New-Object System.Text.UTF8Encoding $false))

Say "Сборка и запуск. Первый раз 15–30 минут: собирается веб-версия и скачиваются модели"
docker compose -f (Join-Path $root "deploy\compose.yaml") -f (Join-Path $root "deploy\tester.yaml") --env-file $envFile up -d --build
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Say "Жду сервер"
for ($i = 0; $i -lt 90; $i++) {
  try {
    Invoke-WebRequest -UseBasicParsing -Uri "http://localhost:8080/healthz" -TimeoutSec 2 | Out-Null
    break
  } catch {
    Start-Sleep -Seconds 2
  }
}

Write-Host ""
Write-Host "Готово." -ForegroundColor Green
Write-Host "  Lodestar:      http://localhost:8080  (с телефона в той же Wi-Fi: http://${ip}:8080)"
Write-Host "  Код вуза:      demo"
Write-Host "  Вход:          раздел «Режим тестировщика» на экране входа"
Write-Host "  Письма:        http://localhost:8025"
if (-not $NoAi) {
  Write-Host "  ИИ-навигатор:  модель $Model ($($values['LODESTAR_LLM_URL'])), эмбеддинги bge-m3"
}
Write-Host "  Обновить или перезапустить: запустите установщик снова"
