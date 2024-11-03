-- +migrate Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE backup_sender_proofs
(
    id                             uuid default uuid_generate_v4() not null primary key,
    enough_balance_proof_body_hash varchar(66) not null,
    last_balance_proof_body        bytea,
    balance_transition_proof_body  bytea,
    created_at                     timestamp with time zone default now() not null
);

-- +migrate Down

DROP TABLE backup_sender_proofs;
