-- +migrate Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE backup_transactions
(
    id               uuid default uuid_generate_v4() not null primary key,
    sender           varchar(255) not null,
    encrypted_tx     text not null,
    block_number     integer not null,
    signature        text not null,
    created_at       timestamp with time zone default now() not null,
    tx_double_hash   text,
    encoding_version integer default 0 not null
);

CREATE INDEX idx_backup_transactions_sender ON backup_transactions (sender);

CREATE INDEX idx_backup_transactions_created_at_block_number ON backup_transactions (created_at, block_number);

CREATE INDEX idx_backup_transactions_created_at_id ON backup_transactions (created_at, id);

-- +migrate Down

DROP TABLE backup_transactions;
