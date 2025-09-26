# SubAggr - Сервис управления подписками для Ed-tech стартапа

REST-сервис для Ed-tech стартапа, решающий задачи агрегации данных об онлайн-подписках пользователей (курсы, тарифы, доп. опции) с возможностью CRUDL-операций и расчёта суммарной стоимости/аналитики.

## Особенности
- Управление подписками (создание, чтение, обновление, удаление, список)
- Расчёт суммарной стоимости подписок за период (финансовая аналитика)
- Фильтрация по пользователю и сервису
- Поддержка UUID для идентификации пользователей
- Swagger-документация API
- Docker-контейнеризация
\- Оптимизировано под Ed-tech предметную область

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

Агрегация (SUMMARY)
POST /SUBS/SUMMARY - Рассчитать сумму подписок за период

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
Расчёт суммы подписок
bash
curl -X POST http://localhost:8000/SUBS/SUMMARY \
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

## Use-cases для Ed-tech
- Ежемесячная финансовая аналитика по пользователю
  - Рассчитать сумму подписок студента за период обучения (например, семестр):
  ```bash
  curl -X POST http://localhost:8000/SUBS/SUMMARY \
    -H "Content-Type: application/json" \
    -d '{
      "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
      "start_date": "09-2025",
      "end_date": "12-2025"
    }'
  ```
- Аналитика по сервису/курсу
  - Посчитать затраты по конкретному сервису (например, платформа вебинаров) для пользователя или всего пула пользователей (omit `user_id`):
  ```bash
  curl -X POST http://localhost:8000/SUBS/SUMMARY \
    -H "Content-Type: application/json" \
    -d '{
      "service_name": "Yandex Plus",
      "start_date": "01-2025",
      "end_date": "06-2025"
    }'
  ```
- Бессрочные подписки и пробные периоды
  - Если `end_date` не задан, подписка считается активной до текущего `end_date` фильтра. Можно сбросить конец подписки через обновление: `{"end_date": ""}`.
- Корректировки тарифа
  - Поддерживается установка цены `0` (например, промо/грант). При создании и обновлении цена должна быть неотрицательной.