-- +migrate Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE backup_transfers
(
    id                   uuid default uuid_generate_v4() not null primary key,
    recipient            varchar(255) not null,
    encrypted_transfer   text not null,
    block_number         integer not null,
    created_at           timestamp with time zone default now() not null,
    transfer_double_hash text
);

CREATE INDEX idx_backup_transfers_recipient ON backup_transfers (recipient);

CREATE INDEX idx_backup_transfers_created_at_block_number ON backup_transfers (created_at, block_number);

CREATE INDEX idx_backup_transfers_created_at_id ON backup_transfers (created_at, id);

-- +migrate Down

DROP TABLE backup_transfers;
