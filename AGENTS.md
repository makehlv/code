# AGENTS.md

## Project overview

CLI-утилита `code` (`github.com/makehlv/code`) для git/workflow: squash коммитов, commit/push по имени ветки, ссылки на Jira и GitLab.

Go 1.22+, stdlib only (нет внешних зависимостей в `go.mod`).

## Layout

```
main.go                 # CLI entrypoint, парсинг команд и флагов
logging.go              # цветной slog.Handler
config/                 # загрузка code_settings.json рядом с бинарником
clients/                # внешние клиенты
  git/                  # обёртка над git CLI (os/exec)
services/
  flow/                 # бизнес-логика команд
```

Слои: `main` → `services` → `clients`. Конфиг и logger прокидываются через конструкторы.

## Commands

| Команда | Назначение |
|---------|------------|
| `squash` | Soft-reset + один коммит; fallback-ветка `code-fallback-*`. Флаги: `--compare` (default `develop`), `--message`, `--push-force` |
| `clean` | Удаляет локальные ветки с префиксом `code-fallback` |
| `commit` | `git add .` + сообщение из имени ветки |
| `push` | При грязном дереве — commit, затем push |
| `jira` | Ссылка Jira из префикса ветки. `--link` — только печать |
| `gitlab` | Ссылка на репозиторий. `--mr` — merge requests, `--link` — только печать |

Формат ветки для сообщений/Jira: `PREFIX/123-description` или `PREFIX-123-description` → `[PREFIX-123] description`.

## Config

Файл `code_settings.json` в каталоге бинарника:

```json
{
  "jiraURL": "https://jira.example.com/browse/",
  "gitlabURL": "https://gitlab.example.com/group/"
}
```

Ключи мапятся на поля `config.Config` (имя поля, lowercase или json-тег). Неизвестный ключ — ошибка.

## Build & run

```bash
go build -o code .
# или
go install github.com/makehlv/code

./code <command> [flags]
```

Тестов пока нет. После добавления: `go test ./...`.

## Conventions

- Ошибки оборачивать через `fmt.Errorf("...: %w", err)` или с контекстом операции.
- Логировать через переданный `*slog.Logger`, не `fmt`/`log` в сервисах (кроме печати ссылок пользователю).
- Git-операции только в `clients/git`; сервисы не вызывают `exec` напрямую (исключение: `OpenInChromeOrPrint` на darwin).
- Новые команды: case в `main.go` + метод на `CodeFlowManageService`.
- Новый внешний клиент: пакет в `clients/`, поле в `clients.Clients`.
- Не коммитить бинарник `code` и секреты; settings лежит рядом с установленным бинарником, не в репо.
