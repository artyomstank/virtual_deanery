# Инструкция по развертыванию и доступу

## Запуск приложения

```bash
docker-compose up -d
```

Приложение запустится на **http://localhost:8080**

## Доступ к фронтенду

### Вариант 1: Локально на той же машине
1. Откройте файл `acl-viewer.html` в браузере
2. В поле "Base URL" должно быть: `http://localhost:8080`
3. Все работает ✅

### Вариант 2: С другой машины или из Docker контейнера

Если открываете фронтенд с другой машины:

1. Узнайте IP адрес машины с приложением:
   ```bash
   # На Linux/Mac:
   hostname -I
   
   # На Windows (PowerShell):
   ipconfig
   ```

2. Замените в `acl-viewer.html` адрес `localhost` на ваш IP:
   ```javascript
   // Было:
   const API_BASE = 'http://localhost:8080';
   
   // Стало (например):
   const API_BASE = 'http://192.168.1.100:8080';
   ```

3. Или используйте поле "Base URL" в фронтенде для смены адреса

## Проверка подключения

```bash
# Проверить, работает ли API
curl http://localhost:8080/health

# Должен вернуть:
# {"status":"ok"}
```

## Логи приложения

```bash
# Смотреть логи API
docker-compose logs -f app

# Смотреть логи базы данных
docker-compose logs -f postgres
```

## Остановка приложения

```bash
docker-compose down
```

## Переменные окружения

Все переменные заданы в `docker-compose.yml`:
- `DB_HOST`: postgres (имя контейнера в сети Docker)
- `DB_PORT`: 5432
- `DB_NAME`: myapp
- `DB_USER`: user
- `DB_PASSWORD`: password
- `JWT_SECRET`: dev-secret-key-change-in-production
- `HTTP_PORT`: 8080

