# Тестирование


```bash

# Запустить Docker контейнеры
docker-compose up -d

# Проверить статус
docker-compose ps
```


---

## Сначала Регистрация администратора

```bash
# Регистрируем пользователя
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "password123"
  }'
```

**Ответ:**
```json
{
  "id": 1,
  "username": "admin",
  "email": "admin@example.com",
  "role": "student",
  "is_active": true,
  "created_at": "2026-05-24T06:18:09Z",
  "updated_at": "2026-05-24T06:18:09Z"
}
```

> Потому что Новый пользователь создается с ролью "student" по умолчанию

---

## Обновление роли на администратора (в БД)

```bash
# Обновляем роль пользователя с ID=1 на админ
docker exec myapp_postgres psql -U user -d myapp -c \
  "UPDATE user_roles SET role_id = (SELECT id FROM roles WHERE name = 'admin') WHERE user_id = 1;"
```

**Ответ:** `UPDATE 1`

---

## Получение JWT токена

```bash
# Получаем токен для администратора
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

**Ответ:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Nzk2OTAyNjQsImlhdCI6MTc3OTYwMzg2NCwidXNlcl9pZCI6MSwicm9sZSI6ImFkbWluIn0.cZzOqPmd7HRkQ8zanZF3wgk3zJg8k5C7acy_hsgRGQg",
  "expires_at": "2026-05-25T06:24:24Z",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin",
    "is_active": true
  }
}
```

**Сохранить токен** из поля `access_token` для последующих запросов


---

## Открыть фронтенд

Откройте файл [acl-viewer.html](acl-viewer.html) в браузере:

1. **Windows:** `file:///c:/uni/isit/Register_LR3/acl-viewer.html`
2. **Вставьте токен** в поле "Токен"
3. **Выбирайте роли** слева - увидите таблицу привилегий

---

## Просмотр всех ролей

```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Nzk2OTAyNjQsImlhdCI6MTc3OTYwMzg2NCwidXNlcl9pZCI6MSwicm9sZSI6ImFkbWluIn0.cZzOqPmd7HRkQ8zanZF3wgk3zJg8k5C7acy_hsgRGQg"

curl -X GET http://localhost:8080/api/v1/admin/roles \
  -H "Authorization: Bearer $TOKEN"
```

**Ответ:**
```json
[
  {
    "id": 1,
    "name": "admin",
    "description": "Administrator with full access"
  },
  {
    "id": 2,
    "name": "teacher",
    "description": "Teacher role"
  },
  {
    "id": 3,
    "name": "student",
    "description": "Student role"
  },
  {
    "id": 4,
    "name": "dean",
    "description": "Dean role"
  }
]
```

---


### Привилегии администратора

```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Nzk2OTAyNjQsImlhdCI6MTc3OTYwMzg2NCwidXNlcl9pZCI6MSwicm9sZSI6ImFkbWluIn0.cZzOqPmd7HRkQ8zanZF3wgk3zJg8k5C7acy_hsgRGQg"

curl -X GET http://localhost:8080/api/v1/admin/acl/admin \
  -H "Authorization: Bearer $TOKEN"
```

### Привилегии преподавателя

```bash
curl -X GET http://localhost:8080/api/v1/admin/acl/teacher \
  -H "Authorization: Bearer $TOKEN"
```

### Привилегии студента

```bash
curl -X GET http://localhost:8080/api/v1/admin/acl/student \
  -H "Authorization: Bearer $TOKEN"
```

### Привилегии декана

```bash
curl -X GET http://localhost:8080/api/v1/admin/acl/dean \
  -H "Authorization: Bearer $TOKEN"
```

**Ответ (пример для student):**
```json
[
  {
    "role_id": 3,
    "resource_id": 3,
    "resource": {
      "id": 3,
      "name": "grades",
      "description": "Grade management"
    },
    "can_read": true,
    "can_write": false,
    "can_delete": false
  },
  {
    "role_id": 3,
    "resource_id": 4,
    "resource": {
      "id": 4,
      "name": "schedule",
      "description": "Schedule management"
    },
    "can_read": true,
    "can_write": false,
    "can_delete": false
  },
  {
    "role_id": 3,
    "resource_id": 5,
    "resource": {
      "id": 5,
      "name": "profile",
      "description": "User profile"
    },
    "can_read": true,
    "can_write": true,
    "can_delete": false
  }
]
```

---

## Полный сценарий в одной команде

```bash
cd c:\uni\isit\Register_LR3 && \
docker-compose up -d && \
sleep 15 && \
curl -s -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","email":"admin@example.com","password":"password123"}' && \
sleep 2 && \
docker exec myapp_postgres psql -U user -d myapp -c \
  "UPDATE user_roles SET role_id = (SELECT id FROM roles WHERE name = 'admin') WHERE user_id = 1;" && \
sleep 2 && \
curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123"}' | python -m json.tool
```

---

## Справка по ролям и привилегиям

### Доступные ресурсы
- `admin` - Админ панель и управление ACL
- `students` - Записи студентов
- `grades` - Управление оценками
- `schedule` - Управление расписанием
- `teachers` - Записи преподавателей
- `profile` - Профиль пользователя
- `reports` - Системные отчеты

### Доступные действия
- `read` (Чтение) - Просмотр ресурса
- `write` (Запись) - Редактирование ресурса
- `delete` (Удаление) - Удаление ресурса

### Матрица привилегий

| Роль | Admin | Students | Grades | Schedule | Teachers | Profile | Reports |
|------|:-----:|:--------:|:------:|:--------:|:--------:|:-------:|:-------:|
| **admin** | RWD | RWD | RWD | RWD | RWD | RWD | RWD |
| **teacher** | - | R | RW | R | - | RW | - |
| **student** | - | - | R | R | - | RW | - |
| **dean** | RW | RW | RW | RW | RW | RW | RW |

> R = Read (Чтение), W = Write (Запись), D = Delete (Удаление), RWD = Все права

---

