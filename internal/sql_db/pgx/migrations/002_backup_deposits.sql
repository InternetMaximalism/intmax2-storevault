-- +migrate Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE backup_deposits
(
    id                  uuid default uuid_generate_v4() not null primary key,
    recipient           varchar(255) not null,
    encrypted_deposit   text not null,
    block_number        integer not null,
    created_at          timestamp with time zone default now() not null,
    deposit_double_hash text
);

CREATE INDEX idx_backup_deposits_recipient ON backup_deposits (recipient);

CREATE INDEX idx_backup_deposits_created_at_block_number ON backup_deposits (created_at, block_number);

CREATE INDEX idx_backup_deposits_created_at_id ON backup_deposits (created_at, id);

-- +migrate Down

DROP TABLE backup_deposits;
