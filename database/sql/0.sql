create table migrations (
    id integer constraint pk_migrations primary key,
    applied_at date not null default now()
);
