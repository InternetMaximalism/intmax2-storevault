package models

import (
	"database/sql"
	"time"
)

type BackupTransaction struct {
	ID              string
	Sender          string
	TxDoubleHash    sql.NullString
	EncryptedTx     string
	EncodingVersion int64
	BlockNumber     int64
	Signature       string
	CreatedAt       time.Time
}

type ListOfBackupTransaction []BackupTransaction
