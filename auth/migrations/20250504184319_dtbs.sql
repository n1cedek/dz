-- +goose Up
create type role as enum ('user','admin');

create table auth(
    id serial primary key,
    name text not null,
    email text not null,
    password text not null,
    role role not null default 'user',
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table if exists auth;
drop type if exists role;

