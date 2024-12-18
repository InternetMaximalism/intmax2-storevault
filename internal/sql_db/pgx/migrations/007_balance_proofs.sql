-- +migrate Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE balance_proofs (
    id uuid not null default uuid_generate_v4(),
    user_state_id uuid not null references user_states(id),
    user_address varchar(66) not null,
    private_state_commitment varchar(66) not null,
    block_number int not null,
    balance_proof bytea not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    PRIMARY KEY (id),
    UNIQUE (user_address, block_number, private_state_commitment),
    UNIQUE (user_state_id)
);

-- +migrate Down

DROP TABLE balance_proofs;
