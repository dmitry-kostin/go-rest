BEGIN;

CREATE TYPE enum_user_role AS ENUM (
    'admin',
    'customer'
    );
ALTER TABLE users
    ADD COLUMN role enum_user_role;

COMMIT;