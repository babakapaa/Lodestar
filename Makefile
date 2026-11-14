.PHONY: up down tester tester-down server-test server-lint app-get app-test app-lint app-web contracts-lint lint test

up:
	docker compose -f deploy/compose.yaml --env-file deploy/.env up -d --build

down:
	docker compose -f deploy/compose.yaml --env-file deploy/.env down

tester:
	docker compose -f deploy/compose.yaml -f deploy/tester.yaml --env-file deploy/.env up -d --build

tester-down:
	docker compose -f deploy/compose.yaml -f deploy/tester.yaml --env-file deploy/.env down

server-test:
	cd server && go test -race ./...

# Тесты базы: LODESTAR_TEST_DATABASE_URL=postgres://user@host:port/postgres make server-test

server-lint:
	cd server && go vet ./... && golangci-lint run

app-get:
	cd app && flutter pub get && flutter gen-l10n

app-test: app-get
	cd app && flutter test

app-lint: app-get
	cd app && dart format --output=none --set-exit-if-changed lib test && flutter analyze

contracts-lint:
	npx --yes @redocly/cli@1.34.3 lint contracts/openapi.yaml

lint: server-lint app-lint contracts-lint

test: server-test app-test

app-web: app-get
	cd app && flutter build web --release --no-web-resources-cdn
