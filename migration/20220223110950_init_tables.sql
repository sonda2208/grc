-- +goose Up
-- +goose StatementBegin
create table accounts (
    id serial primary key,
    created_at timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    deleted_at timestamp without time zone,
    email text not null unique,
    first_name text not null default '',
    last_name text not null default ''
);

create table repositories (
    id serial primary key,
    account_id integer references accounts (id) on delete cascade,
    created_at timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    deleted_at timestamp without time zone,
    name text not null default '',
    url text not null default ''
);

create table scans (
    id serial primary key,
    repo_id integer references repositories (id) on delete cascade,
    created_at timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    scanning_at timestamp without time zone,
    finished_at timestamp without time zone,
    branch text not null default '',
    commit text not null default '',
    status text not null default '',
    message text not null default '',
    findings jsonb not null default '{}'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists scans cascade;
drop table if exists repositories cascade;
drop table if exists accounts cascade;
-- +goose StatementEnd
