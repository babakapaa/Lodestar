# Приложение Lodestar

Одно приложение на Flutter для всех платформ и ролей.

## Папки

| Папка | Что внутри |
| --- | --- |
| `lib/app/` | Сборка приложения: роутер, оболочка с навигацией |
| `lib/core/` | Общее для всех модулей: настройки, навигация |
| `lib/design/` | Дизайн-система: токены четырёх тем, шрифты, компоненты |
| `lib/features/` | Экраны модулей, по папке на модуль |
| `lib/l10n/` | Переводы: казахский, русский, английский |

## Правила

- Цвета, отступы, скругления и длительности анимаций берутся только из `lib/design/tokens/`.
- Строки только из `lib/l10n/app_*.arb`. Новая строка добавляется во все три файла, это проверяет тест.
- Модуль в `features/` не импортирует другой модуль.

## Команды

```sh
flutter pub get
flutter test
tool/web_assets.sh
flutter run -d chrome
```

Нужен Rust из `packages/lodestar_crdt/rust/rust-toolchain.toml`: на нём собирается движок совместного редактирования заметок. Для веба `tool/web_assets.sh` один раз кладёт в `web/` движок в WebAssembly и SQLite для локальной базы, а для `wasm-bindgen` нужна версия из `Cargo.lock`.

Локальная база заметок описана в `lib/features/notes/notes_database.dart`. После изменения таблиц:

```sh
dart run build_runner build --delete-conflicting-outputs
```
