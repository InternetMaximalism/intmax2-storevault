-- +migrate Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_states (
    id uuid not null default uuid_generate_v4(),
    user_address varchar(66) not null,
    encrypted_user_state text not null,
    auth_signature text not null,
    block_number int not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    PRIMARY KEY (id)
);

-- +migrate Down

DROP TABLE user_states;
