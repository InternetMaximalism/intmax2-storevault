-- +migrate Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE backup_balances
(
    id                      uuid default uuid_generate_v4() not null primary key,
    user_address            varchar(66) not null,
    encrypted_balance_proof text not null,
    encrypted_balance_data  text not null,
    encrypted_txs           json not null,
    encrypted_transfers     json not null,
    encrypted_deposits      json not null,
    block_number            integer not null,
    signature               text not null,
    created_at              timestamp with time zone default now() not null
);

CREATE INDEX idx_backup_balances_user_address ON backup_balances (user_address);

-- +migrate Down

DROP TABLE backup_balances;