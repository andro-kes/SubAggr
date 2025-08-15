# SubAggr - Сервис управления подписками

REST-сервис для агрегации данных об онлайн-подписках пользователей с возможностью CRUDL-операций и расчета суммарной стоимости.

## Особенности
- Управление подписками (создание, чтение, обновление, удаление, список)
- Расчет суммарной стоимости подписок за период
- Фильтрация по пользователю и сервису
- Поддержка UUID для идентификации пользователей
- Swagger-документация API
- Docker-контейнеризация

## Требования
- Docker
- Docker Compose

## Запуск приложения

1. Создайте `.env` файл в корне проекта:
```env
POSTGRES_HOST=postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_secure_password
POSTGRES_DB=subaggr
Запустите сервисы:

bash
docker compose up --build
Сервис будет доступен по адресу: http://localhost:8000

API Endpoints
Подписки (SUBS)
POST /SUBS - Создать новую подписку

GET /SUBS - Получить список подписок (с фильтрацией)

GET /SUBS/{id} - Получить подписку по ID

PUT /SUBS/{id} - Обновить подписку

DELETE /SUBS/{id} - Удалить подписку

Агрегация (SUMSUBS)
POST /SUMSUBS - Рассчитать сумму подписок за период

Примеры запросов
Создание подписки
bash
curl -X POST http://localhost:8000/SUBS \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
Расчет суммы подписок
bash
curl -X POST http://localhost:8000/SUMSUBS \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025",
    "end_date": "12-2025"
  }'
Получение списка подписок
bash
curl -X GET "http://localhost:8000/SUBS?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex%20Plus"
Swagger-документация
Документация API доступна после запуска сервиса:
http://localhost:8000/swagger/index.html

Структура проекта
text
├── cmd/                 # Точка входа
├── internal/
|   |── app/             # Запуск приложения
│   ├── database/        # Работа с БД
│   ├── handlers/        # HTTP-обработчики
│   ├── models/          # Модели данных
│   └── utils/           # Вспомогательные утилиты
├── migrations/          # Миграции БД (автоматические)
├── .env.example         # Шаблон конфигурации
├── compose.yaml         # Docker Compose конфигурация
└── Dockerfile           # Конфигурация Docker

Особенности реализации:
- Используется PostgreSQL для хранения данных
- Автоматические миграции при запуске
- Поддержка частичных обновлений
- Валидация UUID и форматов даты
- Логирование операций
- Оптимизированные запросы к БД

Технологический стек:
Go 1.20+
Gin Framework
GORM
PostgreSQL
Docker
Swagger