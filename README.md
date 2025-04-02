# People Credentials Enrichment API

## О проекте
Простой сервис который при создании пользователей пытается угадать их возраст, национальность и пол.
Также поддерживается удаление, редактирование и поиск пользователей.

## Запуск
### Используя Docker (рекомендуется)
1. Установите докер
2. При необходимости отредактируйте порт сервиса в docker-compose.yml
3. Запустите сервис с помощью ```docker-compose up```

По умолчанию сервис работает на http://localhost:8080/

### Запуск нативно
1. Установите go, go-migrate, postgres
2. Отредактируйте энвы системы или укажите их в .env в корне проекта
   
| Поле | Переменная окружения | Значение по умолчанию | Описание |
|------|----------------------|------------------------|----------|
| `ServerPort` | `PEOPLE_CREDENTIALS_SERVER_PORT` | `"8080"` | Порт, на котором запускается HTTP-сервер приложения |
| `DatabasePort` | `PEOPLE_CREDENTIALS_DATABASE_PORT` | `"5432"` | Порт PostgreSQL-сервера |
| `DatabaseUser` | `PEOPLE_CREDENTIALS_DATABASE_USER` | `"postgres"` | Имя пользователя базы данных |
| `DatabasePass` | `PEOPLE_CREDENTIALS_DATABASE_PASSWORD` | `"password"` | Пароль пользователя базы данных |
| `DatabaseName` | `PEOPLE_CREDENTIALS_DATABASE_NAME` | `"user_creds_db"` | Название используемой базы данных |
| `DatabaseHost` | `PEOPLE_CREDENTIALS_DATABASE_HOST` | `"localhost"` | Адрес хоста PostgreSQL |
| `DatabaseSSLMode` | `PEOPLE_CREDENTIALS_DATABASE_SSL_MODE` | `"disable"` | Режим использования SSL при подключении к базе данных |
| `LogLevel` | `PEOPLE_CREDENTIALS_LOG_LEVEL` | `"info"` | Уровень логирования (debug, info, warn, error, fatal) |

3. Создайте пользователя и соответствующую базу данных

4. Запустите миграции, указав актуальные данные вашей базы данных и пользователя
```
migrate -path ./db/migrations -database "postgres://postgres:password@localhost:5432/user_creds_db?sslmode=disable" up 2
```

5. Запустите сервис ``` go run cmd/app/main.go ```

## Примеры использования
### Создание новой записи

**Запрос:**

```http
POST /api/v1/person/create HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
    "name": "vladislav",
    "surname": "bezmaternih",
    "patronymic": "mychailovich"
}
```

**Ответ:**

```http
HTTP/1.1 200 OK
```

---

### Поиск по базе

**Запрос:**

```http
GET /api/v1/search HTTP/1.1
Host: localhost:8080
```

**Ответ:**

```json
[
    {
        "id": 1,
        "name": "vladislav",
        "surname": "bezmaternih",
        "patronymic": "mychailovich",
        "age": 66,
        "gender": "male",
        "nationality": "UA"
    }
]
```

📚 **Полная документация API доступна [здесь](docs/swagger.yaml)**
