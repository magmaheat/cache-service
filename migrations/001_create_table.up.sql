create table if not exists users
(
    id         serial primary key,
    login   varchar(255) not null unique,
    password   varchar(255) not null,
    created_at timestamp    not null default now()
)