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

## Режим тестировщика

Чтобы погонять Lodestar самому, без настройки входа вуза:

```sh
cp deploy/.env.example deploy/.env
make app-web
make tester
```

В `deploy/.env` впишите `LODESTAR_SECRET_KEY` (`openssl rand -base64 32`). `make app-web` нужен только для веб-версии, ему нужны Flutter и Rust, как в разделе «Разработка».

- Веб-версия: `http://localhost:8080`. Код вуза `demo`.
- На экране входа раздел «Режим тестировщика»: вход за студентку, студента, преподавателя или администратора одной кнопкой, без пароля. В профиле тот же раздел меняет пользователя.
- Письма не уходят наружу, их видно в Mailpit: `http://localhost:8025`.
- С телефона: в `deploy/.env` укажите `LODESTAR_PUBLIC_URL=http://<IP компьютера>:8080` и `LODESTAR_S3_PUBLIC_ENDPOINT=http://<IP компьютера>:9000`.
- Звонки: добавьте `--profile calls` и `LIVEKIT_API_SECRET`, `LODESTAR_LIVEKIT_URL` в `deploy/.env`.
- ИИ-навигатор: укажите `LODESTAR_LLM_URL` и `LODESTAR_LLM_MODEL` любого OpenAI-совместимого сервера, например Ollama.

Режим пускает кого угодно без пароля. Не включайте его у настоящего вуза.

## Тестовый вуз

Для показа и приёмки есть тестовый вуз `demo` со своим входом без паролей: студентка, преподаватель и администратор. Его нельзя включать у настоящего вуза.

В `deploy/.env` впишите `DEMO_IDP_CLIENT_SECRET` (любая длинная строка) и `DEMO_IDP_ISSUER`. Это адрес демо-входа, который открывается и с телефона, и с сервера, например `http://192.168.1.10:8090`. `LODESTAR_PUBLIC_URL` должен открываться с телефона так же.

```sh
docker compose -f deploy/compose.yaml --profile demo up -d
docker compose -f deploy/compose.yaml exec -e LODESTAR_OIDC_CLIENT_SECRET=$DEMO_IDP_CLIENT_SECRET server \
  /lodestarctl demo seed -issuer $DEMO_IDP_ISSUER
```

В приложении укажите адрес сервера и код вуза `demo`, нажмите «Войти через Демо-вход» и выберите, кем войти.

В тестовом вузе уже есть курс «Основы биологии» (BIO-101): группа ДЕМО-21 со студенткой, преподаватель, опубликованное задание со сроком через пять дней и пары по понедельникам и средам.

## Разработка

```sh
make server-test
make app-test
make lint
```

Правила работы с кодом описаны в ТЗ, раздел «Правила чистого кода и приёмка», и в `docs/adr/`.
