DROP TABLE IF EXISTS booking_policies;
DROP TABLE IF EXISTS organizations;
DROP TYPE IF EXISTS organization_status;

DROP TRIGGER IF EXISTS trg_create_default_booking_policy ON organizations;
DROP FUNCTION IF EXISTS create_default_booking_policy();