create table migrations (
    id integer constraint pk_migrations primary key,
    applied_at date not null default now()
);

create table games (
    id uuid constraint pk_game primary key default gen_random_uuid(),
    title varchar(200) not null,
    title_normalised varchar(200) generated always as ( upper(title) ) stored,
    description varchar(2000) null,
    archived bool not null
);

create table platforms (
    id uuid constraint pk_platform primary key default gen_random_uuid(),
    name varchar(200) not null constraint ix_platform_name unique,
    name_normalised varchar(200) generated always as ( upper(name) ) stored,
    short_name varchar(10) not null constraint ix_platform_short_name unique,
    short_name_normalised varchar(10) generated always as ( upper(short_name) ) stored
);

create table game_releases (
    id uuid constraint pk_game_releases primary key default gen_random_uuid(),
    game_id uuid not null references games(id) on delete cascade,
    title_override varchar(200) null,
    title_override_normalised varchar(200) generated always as ( upper(title_override) ) stored,
    description varchar(2000) null,
    release_date date null,
    release_date_unknown bool not null
);

create table game_release_platforms (
    platform_id uuid not null references platforms(id),
    game_release_id uuid not null references game_releases(id),
    constraint ix_game_release_platforms_unique unique (platform_id, game_release_id)
);