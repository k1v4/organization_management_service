# organization_management_service

Сервис управления организациями и политиками бронирования.

## Запуск

### Переменные окружения

Создай файл `.env` в корне проекта:

```env
KEYCLOAK_ISSUER=https://keycloak.vts-platform.ru/realms/organization-bookings
MEMBERSHIP_SERVICE_URL=http://membership-service:8080
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
DB_URL=postgres://user:password@localhost:5432/orgservice
```

### Запуск сервиса

```bash
docker compose up -d
```

Миграции применятся при запуске

---

## API

Базовый URL: `/api/v1`

Все запросы требуют заголовок:
```
Authorization: Bearer <JWT токен>
```

### GET /api/organizations/health

есть отличающийся эндпоинт, который отвечает за проверку того, что сервис встал

---

### Организации

#### POST /api/v1/organizations

Создать организацию. Инициатор автоматически становится владельцем с ролью `ORG_OWNER`.

**Тело запроса:**
```json
{
    "name": "Рога и копыта",
    "description": "Описание организации"
}
```

**Ответ** `201 Created`:
```json
{
    "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "name": "Рога и копыта",
    "description": "Описание организации",
    "status": "active",
    "owner_identity_id": "auth0|abc123",
    "created_at": "2026-03-21T10:00:00Z",
    "updated_at": "2026-03-21T10:00:00Z"
}
```

---

#### GET /api/v1/organizations/:orgId

Получить организацию по идентификатору.

**Параметры пути:**
- `orgId` — UUID организации

**Ответ** `200 OK`:
```json
{
    "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "name": "Рога и копыта",
    "description": "Описание организации",
    "status": "active",
    "owner_identity_id": "auth0|abc123",
    "created_at": "2026-03-21T10:00:00Z",
    "updated_at": "2026-03-21T10:00:00Z"
}
```

---

#### PUT /api/v1/organizations/:orgId

Обновить название и описание организации. Требует право `ORG_UPDATE`.

**Параметры пути:**
- `orgId` — UUID организации

**Тело запроса:**
```json
{
    "name": "Новое название",
    "description": "Новое описание"
}
```

**Ответ** `200 OK`:
```json
{
    "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "name": "Новое название",
    "description": "Новое описание",
    "status": "active",
    "owner_identity_id": "auth0|abc123",
    "created_at": "2026-03-21T10:00:00Z",
    "updated_at": "2026-03-21T11:00:00Z"
}
```

---

#### DELETE /api/v1/organizations/:orgId

Архивировать организацию (мягкое удаление). Требует право `ORG_DELETE`.

**Параметры пути:**
- `orgId` — UUID организации

**Ответ** `204 No Content`

---

#### PUT /api/v1/organizations/:orgId/owner

Сменить владельца организации. Требует право `ORG_OWNER_CHANGE`.

**Параметры пути:**
- `orgId` — UUID организации

**Тело запроса:**
```json
{
    "new_owner_identity_id": "auth0|xyz789"
}
```

**Ответ** `200 OK`:
```json
{
    "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "name": "Рога и копыта",
    "description": "Описание организации",
    "status": "active",
    "owner_identity_id": "auth0|xyz789",
    "created_at": "2026-03-21T10:00:00Z",
    "updated_at": "2026-03-21T12:00:00Z"
}
```

---

### Политики бронирования

Политика создаётся автоматически при создании организации с дефолтными значениями. Отдельного эндпоинта для создания нет.

**Дефолтные значения:**
- `max_booking_duration_min` — 480 (8 часов)
- `booking_window_days` — 30 дней
- `max_active_bookings_per_user` — 5

---

#### GET /api/v1/organizations/:orgId/policy

Получить политику бронирования организации. Требует право `POLICIES_LIST`.

**Параметры пути:**
- `orgId` — UUID организации

**Ответ** `200 OK`:
```json
{
    "id": 1,
    "organization_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "max_booking_duration_min": 480,
    "booking_window_days": 30,
    "max_active_bookings_per_user": 5,
    "created_at": "2026-03-21T10:00:00Z",
    "updated_at": "2026-03-21T10:00:00Z"
}
```

---

#### PUT /api/v1/organizations/:orgId/policy

Обновить политику бронирования. Требует право `POLICIES_MANAGE`. Все поля опциональны — передавай только те, которые нужно изменить.

**Параметры пути:**
- `orgId` — UUID организации

**Тело запроса:**
```json
{
    "max_booking_duration_min": 120,
    "booking_window_days": 14,
    "max_active_bookings_per_user": 3
}
```

**Ответ** `200 OK`:
```json
{
    "id": 1,
    "organization_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "max_booking_duration_min": 120,
    "booking_window_days": 14,
    "max_active_bookings_per_user": 3,
    "created_at": "2026-03-21T10:00:00Z",
    "updated_at": "2026-03-21T13:00:00Z"
}
```

---

## Коды ошибок

| Код | Описание |
|-----|----------|
| `400` | Невалидное тело запроса или параметры |
| `401` | Токен отсутствует или невалиден |
| `403` | Недостаточно прав |
| `404` | Организация не найдена |
| `500` | Внутренняя ошибка сервера |

**Формат ошибки:**
```json
{
    "error": "описание ошибки"
}
```

---

## Схема БД

```
organizations
├── id                UUID PK
├── name              VARCHAR(255)
├── description       TEXT
├── status            ENUM(active, archived)
├── owner_identity_id VARCHAR(255)
├── created_at        TIMESTAMPTZ
└── updated_at        TIMESTAMPTZ

booking_policies
├── id                          SERIAL PK
├── organization_id             UUID FK → organizations.id
├── max_booking_duration_min    INT (default: 480)
├── booking_window_days         INT (default: 30)
├── max_active_bookings_per_user INT (default: 5)
├── created_at                  TIMESTAMPTZ
└── updated_at                  TIMESTAMPTZ
```

---

## Интеграции

Сервис взаимодействует с **OrgMembershipService** для проверки прав и управления членством.

| Операция | Вызов |
|----------|-------|
| Проверка прав | `POST /api/internal/authorization/check` |
| Проверка пользователя | `GET /api/internal/users/by-identity/{identityId}` |
| Назначение владельца | `POST /api/internal/organizations/owner` |
| Получение членства | `GET /api/internal/organizations/{orgId}/users/{identityId}/membership` |
| Снятие роли | `DELETE /api/Organizations/{orgId}/members/{membershipId}/roles/{roleCode}` |