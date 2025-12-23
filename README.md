# Go Chat - CLI-чат c WebSocket

Проект чата на Go с использованием WebSocket сервера для real-time обмена сообщениями.

## Архитектура

Проект следует принципам чистой архитектуры с использованием `internal` пакета для приватных компонентов:

```
gochat/
├── internal/        # Приватные компоненты приложения
│   ├── domain/      # Доменные сущности и интерфейсы репозиториев
│   ├── repository/ # Реализации репозиториев (in-memory)
│   ├── usecase/    # Бизнес-логика
│   └── delivery/   # HTTP handlers и WebSocket сервер
├── cmd/server/      # Точка входа сервера
├── client/          # CLI клиент
└── Makefile         # Команды для запуска и сборки
```

## Требования

- Go 1.21 или выше

## Быстрый старт

```bash
# 1. Запустить сервер (в первом терминале)
make server

# 2. Запустить клиент (во втором терминале)
make client
```

## Настройка

### Переменные окружения

Проект поддерживает загрузку переменных окружения из `.env` файла (опционально).

Создайте `.env` файл в корне проекта:

```env
# Server Configuration
PORT=8080
HOST=0.0.0.0

# Client Configuration
SERVER_URL=http://localhost:8080
WS_URL=ws://localhost:8080/ws
```

#### Для сервера:
- `PORT` - Порт сервера (по умолчанию: 8080)
- `HOST` - Хост для прослушивания (по умолчанию: 0.0.0.0 - все интерфейсы)

#### Для клиента:
- `SERVER_URL` - Адрес сервера (по умолчанию: http://localhost:8080)
- `WS_URL` - WebSocket URL (по умолчанию: автоматически формируется из SERVER_URL)

### Запуск с параметрами

Переменные окружения можно задавать через `.env` файл или через командную строку:

```bash
# Сервер (Linux/Mac)
PORT=3000 HOST=0.0.0.0 make server

# Сервер (Windows PowerShell)
$env:PORT="3000"; $env:HOST="0.0.0.0"; make server

# Клиент (подключение к удаленному серверу, Linux/Mac)
SERVER_URL=http://192.168.1.100:8080 make client

# Клиент (Windows PowerShell)
$env:SERVER_URL="http://192.168.1.100:8080"; make client
```

Приоритет: переменные из командной строки > переменные из `.env` > значения по умолчанию

## Команды Makefile

```bash
make server  # Запустить сервер
make client  # Запустить клиент
make test    # Запустить тесты
```

## API Endpoints

### Пользователи

- `POST /api/users/register` - Регистрация пользователя
  ```json
  {
    "username": "john_doe"
  }
  ```

- `GET /api/users/get?id={user_id}` - Получение пользователя по ID

### Комнаты

- `POST /api/rooms/create` - Создание комнаты
  ```json
  {
    "name": "General"
  }
  ```

- `GET /api/rooms/get?id={room_id}` - Получение комнаты по ID
- `GET /api/rooms/all` - Получение всех комнат

### Сообщения

- `POST /api/messages/send?room_id={room_id}&user_id={user_id}` - Отправка сообщения
  ```json
  {
    "content": "Hello, world!"
  }
  ```

- `GET /api/messages/history?room_id={room_id}&limit=50&offset=0` - Получение истории сообщений

### WebSocket

- `GET /ws?room_id={room_id}&user_id={user_id}` - Подключение к WebSocket для real-time сообщений

## Использование по сети

Сервер по умолчанию слушает на всех интерфейсах (`0.0.0.0`), что позволяет подключаться с других компьютеров в сети.

## Команды клиента

После запуска клиента доступны следующие команды:

- `/rooms` - Показать все комнаты
- `/create <name>` - Создать новую комнату
- `/join <number>` - Присоединиться к комнате по номеру
- `/leave` - Покинуть текущую комнату
- `/history [limit]` - Показать историю сообщений (по умолчанию: 10)
- `/help` - Показать справку
- `/exit` - Выйти из приложения

## Тестирование

```bash
# Запустить все тесты
make test
```

## Лицензия

MIT
