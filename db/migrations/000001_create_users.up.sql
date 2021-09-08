CREATE TABLE IF NOT EXISTS users
(
    id         serial primary key,
    email      varchar(300) unique not null,
    first_name varchar(50)         not null,
    last_name  varchar(50),
    created_at timestamp default now()
);