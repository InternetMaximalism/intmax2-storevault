package models

import (
	"database/sql"
	"time"
)

type BackupDeposit struct {
	ID                string
	Recipient         string
	DepositDoubleHash sql.NullString
	EncryptedDeposit  string
	BlockNumber       int64
	CreatedAt         time.Time
}

type ListOfBackupDeposit []BackupDeposit
