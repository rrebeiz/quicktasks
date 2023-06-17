create table if not exists tasks (
    id bigserial primary key,
    title text not null,
    description text not null,
    completed boolean not null,
    version integer not null default 1,
    created_at timestamp(0) without time zone not null default now(),
    updated_at timestamp(0) without time zone not null default now()
);