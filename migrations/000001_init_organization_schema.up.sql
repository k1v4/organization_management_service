CREATE TYPE organization_status AS ENUM ('active', 'archived');

CREATE TABLE organizations (
                               id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               name                VARCHAR(255) NOT NULL,
                               description         TEXT,
                               status              organization_status NOT NULL DEFAULT 'active',
                               owner_identity_id   VARCHAR(255) NOT NULL,
                               created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
                               updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE booking_policies (
                                  id                           SERIAL PRIMARY KEY,
                                  organization_id              UUID NOT NULL UNIQUE REFERENCES organizations(id) ON DELETE CASCADE,
    -- максимальная длительность одного бронирования в минутах (default: 8 часов)
                                  max_booking_duration_min     INT NOT NULL DEFAULT 480,
    -- на сколько дней вперёд можно создавать бронирования (default: 30 дней)
                                  booking_window_days          INT NOT NULL DEFAULT 30,
    -- максимальное количество активных бронирований на одного пользователя одновременно
                                  max_active_bookings_per_user INT NOT NULL DEFAULT 5,
                                  created_at                   TIMESTAMPTZ NOT NULL DEFAULT now(),
                                  updated_at                   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_organizations_status ON organizations(status);
CREATE INDEX idx_organizations_owner  ON organizations(owner_identity_id);