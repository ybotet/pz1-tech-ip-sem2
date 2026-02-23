# Документация API для практического занятия №1

## Общая информация

### Базовые URL

- **Auth Service:** `http://localhost:8081`
- **Tasks Service:** `http://localhost:8082`

### Форматы данных

- Все запросы и ответы передаются в формате JSON
- Кодировка: UTF-8
- Content-Type: `application/json`

### Заголовки общие

| Заголовок       | Описание                       | Обязательность                  |
| --------------- | ------------------------------ | ------------------------------- |
| `Content-Type`  | Должен быть `application/json` | Для POST/PATCH                  |
| `X-Request-ID`  | UUID для трассировки запроса   | Опционально                     |
| `Authorization` | `Bearer <token>`               | Для защищенных эндпоинтов Tasks |

### Коды ответов HTTP

| Код                         | Описание                                              |
| --------------------------- | ----------------------------------------------------- |
| `200 OK`                    | Успешный запрос (GET, PATCH)                          |
| `201 Created`               | Ресурс успешно создан (POST)                          |
| `204 No Content`            | Успешный запрос без тела ответа (DELETE)              |
| `400 Bad Request`           | Неверный формат запроса                               |
| `401 Unauthorized`          | Отсутствует или невалидный токен                      |
| `403 Forbidden`             | Недостаточно прав (токен валиден, но доступ запрещен) |
| `404 Not Found`             | Ресурс не найден                                      |
| `500 Internal Server Error` | Внутренняя ошибка сервера                             |

---

## Auth Service (порт: 8081)

Сервис отвечает за аутентификацию и выдачу токенов доступа.

---

### 1. Вход в систему (Login)

Получение токена доступа по упрощенным учетным данным.

**Эндпоинт:** `POST /v1/auth/login`

#### Запрос

**Заголовки:**

```
Content-Type: application/json
X-Request-ID: req-001
```

**Тело запроса:**

```json
{
  "username": "student",
  "password": "student"
}
```

| Поле       | Тип    | Описание                                          |
| ---------- | ------ | ------------------------------------------------- |
| `username` | string | Имя пользователя (в демо-версии только "student") |
| `password` | string | Пароль (в демо-версии только "student")           |

#### Успешный ответ (200 OK)

```json
{
  "access_token": "demo-token",
  "token_type": "Bearer"
}
```

| Поле           | Тип    | Описание                               |
| -------------- | ------ | -------------------------------------- |
| `access_token` | string | Токен доступа для последующих запросов |
| `token_type`   | string | Тип токена (всегда "Bearer")           |

#### Ответы об ошибках

**400 Bad Request**

```json
{
  "error": "invalid request"
}
```

**401 Unauthorized**

```json
{
  "error": "invalid credentials"
}
```

#### Примеры curl

**Успешный запрос:**

```bash
curl -s -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: req-001" \
  -d '{"username":"student","password":"student"}' | jq
```

**Неверные учетные данные:**

```bash
curl -i -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"wrong"}'
```

---

### 2. Проверка токена (Verify)

Проверяет валидность переданного токена доступа.

**Эндпоинт:** `GET /v1/auth/verify`

#### Запрос

**Заголовки:**

```
Authorization: Bearer demo-token
X-Request-ID: req-002
```

| Заголовок       | Описание                         |
| --------------- | -------------------------------- |
| `Authorization` | Токен в формате `Bearer <token>` |
| `X-Request-ID`  | ID для трассировки (опционально) |

#### Успешный ответ (200 OK)

```json
{
  "valid": true,
  "subject": "student"
}
```

| Поле      | Тип     | Описание                   |
| --------- | ------- | -------------------------- |
| `valid`   | boolean | Результат проверки токена  |
| `subject` | string  | Идентификатор пользователя |

#### Ответ об ошибке (401 Unauthorized)

```json
{
  "valid": false,
  "error": "unauthorized"
}
```

Или с конкретной причиной:

```json
{
  "valid": false,
  "error": "missing token"
}
```

```json
{
  "valid": false,
  "error": "invalid auth header"
}
```

```json
{
  "valid": false,
  "error": "invalid token"
}
```

#### Примеры curl

**Проверка валидного токена:**

```bash
curl -i http://localhost:8081/v1/auth/verify \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-002"
```

**Проверка без токена:**

```bash
curl -i http://localhost:8081/v1/auth/verify
```

**Проверка с невалидным токеном:**

```bash
curl -i http://localhost:8081/v1/auth/verify \
  -H "Authorization: Bearer wrong-token"
```

---

## Tasks Service (порт: 8082)

Сервис управления задачами. Все защищенные эндпоинты требуют валидный токен.

### Общие требования для защищенных эндпоинтов

**Обязательные заголовки:**

```
Authorization: Bearer demo-token
X-Request-ID: req-003
```

**Возможные ошибки авторизации:**

**401 Unauthorized**

```json
{
  "error": "missing authorization header"
}
```

```json
{
  "error": "invalid authorization header format"
}
```

```json
{
  "error": "invalid token"
}
```

**500 Internal Server Error** (при недоступности Auth Service)

```json
{
  "error": "authorization service unavailable"
}
```

---

### 1. Создание задачи

**Эндпоинт:** `POST /v1/tasks`

#### Запрос

**Заголовки:**

```
Authorization: Bearer demo-token
Content-Type: application/json
X-Request-ID: req-003
```

**Тело запроса:**

```json
{
  "title": "Изучить материал",
  "description": "Прочитать лекцию по микросервисам",
  "due_date": "2026-01-15"
}
```

| Поле          | Тип    | Обязательность | Описание                             |
| ------------- | ------ | -------------- | ------------------------------------ |
| `title`       | string | Да             | Заголовок задачи                     |
| `description` | string | Нет            | Подробное описание                   |
| `due_date`    | string | Нет            | Дата выполнения (формат: YYYY-MM-DD) |

#### Успешный ответ (201 Created)

```json
{
  "id": "t_001",
  "title": "Изучить материал",
  "description": "Прочитать лекцию по микросервисам",
  "due_date": "2026-01-15",
  "done": false
}
```

| Поле          | Тип     | Описание                        |
| ------------- | ------- | ------------------------------- |
| `id`          | string  | Уникальный идентификатор задачи |
| `title`       | string  | Заголовок задачи                |
| `description` | string  | Описание задачи                 |
| `due_date`    | string  | Дата выполнения                 |
| `done`        | boolean | Статус выполнения               |

#### Ошибки

**400 Bad Request**

```json
{
  "error": "invalid request"
}
```

#### Пример curl

```bash
curl -i -X POST http://localhost:8082/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-003" \
  -d '{
    "title": "Изучить материал",
    "description": "Прочитать лекцию по микросервисам",
    "due_date": "2026-01-15"
  }'
```

---

### 2. Получение списка задач

**Эндпоинт:** `GET /v1/tasks`

#### Запрос

**Заголовки:**

```
Authorization: Bearer demo-token
X-Request-ID: req-004
```

#### Успешный ответ (200 OK)

```json
[
  {
    "id": "t_001",
    "title": "Изучить материал",
    "done": false
  },
  {
    "id": "t_002",
    "title": "Выполнить практику",
    "done": true
  }
]
```

#### Пример curl

```bash
curl -i http://localhost:8082/v1/tasks \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-004"
```

---

### 3. Получение задачи по ID

**Эндпоинт:** `GET /v1/tasks/{id}`

#### Запрос

**Заголовки:**

```
Authorization: Bearer demo-token
X-Request-ID: req-005
```

**Параметры пути:**
| Параметр | Описание |
|----------|----------|
| `id` | Идентификатор задачи (например, `t_001`) |

#### Успешный ответ (200 OK)

```json
{
  "id": "t_001",
  "title": "Изучить материал",
  "description": "Прочитать лекцию по микросервисам",
  "due_date": "2026-01-15",
  "done": false
}
```

#### Ошибки

**404 Not Found**

```json
{
  "error": "task not found"
}
```

#### Пример curl

```bash
curl -i http://localhost:8082/v1/tasks/t_001 \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-005"
```

---

### 4. Обновление задачи

**Эндпоинт:** `PATCH /v1/tasks/{id}`

#### Запрос

**Заголовки:**

```
Authorization: Bearer demo-token
Content-Type: application/json
X-Request-ID: req-006
```

**Параметры пути:**
| Параметр | Описание |
|----------|----------|
| `id` | Идентификатор задачи (например, `t_001`) |

**Тело запроса:** (все поля опциональны)

```json
{
  "title": "Изучить материал (обновлено)",
  "description": "Добавить конспект",
  "due_date": "2026-01-20",
  "done": true
}
```

| Поле          | Тип     | Описание                |
| ------------- | ------- | ----------------------- |
| `title`       | string  | Новый заголовок         |
| `description` | string  | Новое описание          |
| `due_date`    | string  | Новая дата выполнения   |
| `done`        | boolean | Новый статус выполнения |

#### Успешный ответ (200 OK)

Возвращает обновленную задачу:

```json
{
  "id": "t_001",
  "title": "Изучить материал (обновлено)",
  "description": "Добавить конспект",
  "due_date": "2026-01-20",
  "done": true
}
```

#### Ошибки

**400 Bad Request**

```json
{
  "error": "invalid request"
}
```

**404 Not Found**

```json
{
  "error": "task not found"
}
```

#### Пример curl

```bash
curl -i -X PATCH http://localhost:8082/v1/tasks/t_001 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-006" \
  -d '{
    "title": "Изучить материал (обновлено)",
    "done": true
  }'
```

---

### 5. Удаление задачи

**Эндпоинт:** `DELETE /v1/tasks/{id}`

#### Запрос

**Заголовки:**

```
Authorization: Bearer demo-token
X-Request-ID: req-007
```

**Параметры пути:**
| Параметр | Описание |
|----------|----------|
| `id` | Идентификатор задачи (например, `t_001`) |

#### Успешный ответ (204 No Content)

Ответ без тела.

#### Ошибки

**404 Not Found**

```json
{
  "error": "task not found"
}
```

#### Пример curl

```bash
curl -i -X DELETE http://localhost:8082/v1/tasks/t_001 \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-007"
```

---

## Полный сценарий тестирования

```bash
#!/bin/bash

echo "Тестирование микросервисов"
echo "============================="

# 1. Получение токена
echo -e "\n 1. Получение токена:"
TOKEN=$(curl -s -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: test-001" \
  -d '{"username":"student","password":"student"}' | jq -r '.access_token')
echo "   Токен: $TOKEN"

# 2. Проверка токена
echo -e "\n 2. Проверка токена:"
curl -s -X GET http://localhost:8081/v1/auth/verify \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Request-ID: test-002" | jq

# 3. Создание задачи
echo -e "\n 3. Создание задачи:"
TASK_ID=$(curl -s -X POST http://localhost:8082/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Request-ID: test-003" \
  -d '{"title":"Тестовая задача","description":"Описание"}' | jq -r '.id')
echo "   ID задачи: $TASK_ID"

# 4. Получение списка задач
echo -e "\n 4. Список задач:"
curl -s -X GET http://localhost:8082/v1/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Request-ID: test-004" | jq

# 5. Получение задачи по ID
echo -e "\n 5. Задача $TASK_ID:"
curl -s -X GET http://localhost:8082/v1/tasks/$TASK_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Request-ID: test-005" | jq

# 6. Обновление задачи
echo -e "\n 6. Обновление задачи:"
curl -s -X PATCH http://localhost:8082/v1/tasks/$TASK_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Request-ID: test-006" \
  -d '{"done":true}' | jq

# 7. Удаление задачи
echo -e "\n 7. Удаление задачи:"
curl -i -X DELETE http://localhost:8082/v1/tasks/$TASK_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Request-ID: test-007"

echo -e "\n Тестирование завершено!"
```

---

## Сводная таблица эндпоинтов

### Auth Service

| Метод | Путь              | Описание         | Тело запроса              | Успешный ответ           |
| ----- | ----------------- | ---------------- | ------------------------- | ------------------------ |
| POST  | `/v1/auth/login`  | Получение токена | `{"username","password"}` | `200` + токен            |
| GET   | `/v1/auth/verify` | Проверка токена  | -                         | `200` + `{"valid":true}` |

### Tasks Service

| Метод  | Путь             | Описание          | Тело запроса                         | Успешный ответ       |
| ------ | ---------------- | ----------------- | ------------------------------------ | -------------------- |
| POST   | `/v1/tasks`      | Создание задачи   | `{"title","description","due_date"}` | `201` + задача       |
| GET    | `/v1/tasks`      | Список задач      | -                                    | `200` + массив задач |
| GET    | `/v1/tasks/{id}` | Задача по ID      | -                                    | `200` + задача       |
| PATCH  | `/v1/tasks/{id}` | Обновление задачи | Любые поля задачи                    | `200` + задача       |
| DELETE | `/v1/tasks/{id}` | Удаление задачи   | -                                    | `204`                |

---

## Переменные окружения

### Auth Service

| Переменная  | Значение по умолчанию | Описание              |
| ----------- | --------------------- | --------------------- |
| `AUTH_PORT` | `8081`                | Порт для HTTP сервера |

### Tasks Service

| Переменная      | Значение по умолчанию   | Описание                 |
| --------------- | ----------------------- | ------------------------ |
| `TASKS_PORT`    | `8082`                  | Порт для HTTP сервера    |
| `AUTH_BASE_URL` | `http://localhost:8081` | Базовый URL Auth сервиса |

---

## Примеры с ошибками

### 1. Запрос без токена

```bash
curl -i http://localhost:8082/v1/tasks
```

**Ответ:**

```
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{"error":"missing authorization header"}
```

### 2. Запрос с несуществующим токеном

```bash
curl -i http://localhost:8082/v1/tasks \
  -H "Authorization: Bearer wrong-token"
```

**Ответ:**

```
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{"error":"invalid token"}
```

### 3. Запрос к несуществующей задаче

```bash
curl -i http://localhost:8082/v1/tasks/t_999 \
  -H "Authorization: Bearer demo-token"
```

**Ответ:**

```
HTTP/1.1 404 Not Found
Content-Type: application/json

{"error":"task not found"}
```

---

## Примечания

1. **Таймауты:** Tasks Service ожидает ответ от Auth Service не более 2-3 секунд
2. **Request-ID:** Рекомендуется всегда передавать для трассировки
3. **Хранение данных:** В демо-версии данные хранятся в памяти (теряются при перезапуске)
4. **Токены:** В демо-версии используется фиксированный токен "demo-token"
