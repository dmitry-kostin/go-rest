CREATE TYPE enum_user_role AS ENUM ('admin','customer');

CREATE TABLE IF NOT EXISTS users
(
    id          uuid                not null primary key,
    identity_id uuid                not null,
    email       varchar(255) unique not null,
    role        enum_user_role      not null,
    first_name  varchar(50)         not null,
    last_name   varchar(50),
    created_at  timestamp default now(),
    updated_at  timestamp default now()
);