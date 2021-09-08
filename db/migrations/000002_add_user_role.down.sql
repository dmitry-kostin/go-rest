BEGIN;

ALTER TABLE users DROP COLUMN role;
DROP TYPE enum_user_role;

COMMIT;