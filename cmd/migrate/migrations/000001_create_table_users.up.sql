create table if not exists users(
    id serial primary key,
    username varchar(255) not null unique,
    fullname varchar(255) not null,
    email varchar(255) not null unique,
    password bytea not null,
    is_active boolean not null default false,
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now()
);