-- +goose Up
-- +goose StatementBegin
create table users (
    email text primary key,
    id text not null,
    verificated bool default false,
    baned bool default false,
    ssh_sign text not null,
    deployer text not null,
    create_dt timestamptz default now()
);

create table deployment (
    id text primary key,
    deployer text not null,
    name text not null,
    state text not null,
    ip text not null,
    last_connection timestamptz not null,
    create_dt timestamptz default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE deployment;
-- +goose StatementEnd