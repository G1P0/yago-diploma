# TODO-планировщик

Мини-приложение на Go: хранит задачи в SQLite, поддерживает повторения (d, y, w, m), редактирование, выполнение и удаление. Есть авторизация через пароль (JWT в куках). Веб-интерфейс лежит в директории `web`.

## ⭐ Выполнено со звёздочкой
- Переменные окружения: `TODO_PORT`, `TODO_DBFILE`, `TODO_PASSWORD`
- Полная поддержка правил повторений (w, m)
- JWT-аутентификация + middleware
- Dockerfile + запуск контейнера

## Локальный запуск
```bash
go mod tidy

# Примеры настроек
export TODO_PASSWORD=12345        # пароль для входа (пусто = аутентификация отключена)
export TODO_PORT=7540             # порт сервера, по умолчанию 7540
export TODO_DBFILE=./scheduler.db # путь к базе данных, по умолчанию ./scheduler.db

go run ./cmd/main.go
```
Будет доступен по адресу: `http://localhost:7540/`

## Тесты

Настройки: `tests/settings.go`

- `Port` - порт сервера (по умолчанию `7540`)
- `DBFile` - путь к базе (`../scheduler.db`)
- `FullNextDate` = `true` - проверка всех правил повторов (реализовано)
- `Search` = `true` - проверка поиска (реализовано)
- `Token` - JWT токен, если включена аутентификация (реализовано)

Если аутентификация выключена (`TODO_PASSWORD` не задан):

```bash
go test ./tests
```

Если аутентификация включена:

1. Получить токен через `/api/signin`
2. Вставить его в `tests/settings.go`:

```bash
var Token = <токен>
```

3. Запустить:

```bash
go test ./tests
```

## Docker

Сборка:

```bash
docker build -t todo-app .
```

Запуск:

```bash
docker run --rm -d --name todo-app -e TODO_PASSWORD=12345 -p 7540:7540 todo-app
```

## Структура

```bash
internal/server/   - запуск HTTP-сервера, маршрутизация, статика
internal/handlers/ - отдача HTML/JS/CSS файлов
pkg/api/           - API-хендлеры (/api/task, /api/tasks, /api/done, /api/nextdate, /api/signin)
pkg/db/            - работа с SQLite: инициализация, CRUD-операции, поиск
pkg/auth/          - генерация и проверка JWT-токенов
web/               - статические файлы фронтенда
tests/             - автотесты итогового задания
Dockerfile         - контейнеризация приложения

```