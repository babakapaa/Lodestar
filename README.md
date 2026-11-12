# Lodestar

ИИ-кампус для вуза: ИИ-навигатор ведёт студента к знаниям лестницей подсказок, а не пишет работу за него.

## Устройство

| Папка | Что внутри |
| --- | --- |
| `app/` | Приложение на Flutter для всех платформ и ролей |
| `server/` | Сервер на Go: ядро и модули в одной программе |
| `contracts/` | Контракт API (OpenAPI) и схемы событий |
| `deploy/` | Установка: Docker Compose для пилота |
| `docs/adr/` | Записи архитектурных решений |

## Быстрый старт

```sh
cp deploy/.env.example deploy/.env
make up
```

В `deploy/.env` впишите `LODESTAR_SECRET_KEY`, его выдаёт `openssl rand -base64 32`.

Сервер поднимается на `http://localhost:8080`, проверка: `curl localhost:8080/healthz`.

## Подключить вуз

```sh
docker compose -f deploy/compose.yaml exec server /lodestarctl tenant create -slug kaznu -name "КазНУ"
docker compose -f deploy/compose.yaml exec -e LODESTAR_OIDC_CLIENT_SECRET=... server \
  /lodestarctl oidc set -tenant kaznu -key microsoft -name "Microsoft" -issuer https://login.microsoftonline.com/<id>/v2.0 -client-id <id>
```

Команда напечатает адрес возврата, его нужно указать в настройках приложения у провайдера. Когда администратор вуза войдёт первый раз, выдайте ему высшую роль:

```sh
docker compose -f deploy/compose.yaml exec server /lodestarctl role grant-top -tenant kaznu -email admin@kaznu.kz
```

## Тестовый вуз

Для показа и приёмки есть тестовый вуз `demo` со своим входом без паролей: студентка, преподаватель и администратор. Его нельзя включать у настоящего вуза.

В `deploy/.env` впишите `DEMO_IDP_CLIENT_SECRET` (любая длинная строка) и `DEMO_IDP_ISSUER`. Это адрес демо-входа, который открывается и с телефона, и с сервера, например `http://192.168.1.10:8090`. `LODESTAR_PUBLIC_URL` должен открываться с телефона так же.

```sh
docker compose -f deploy/compose.yaml --profile demo up -d
docker compose -f deploy/compose.yaml exec -e LODESTAR_OIDC_CLIENT_SECRET=$DEMO_IDP_CLIENT_SECRET server \
  /lodestarctl demo seed -issuer $DEMO_IDP_ISSUER
```

В приложении укажите адрес сервера и код вуза `demo`, нажмите «Войти через Демо-вход» и выберите, кем войти.

## Разработка

```sh
make server-test
make app-test
make lint
```

Правила работы с кодом описаны в ТЗ, раздел «Правила чистого кода и приёмка», и в `docs/adr/`.
