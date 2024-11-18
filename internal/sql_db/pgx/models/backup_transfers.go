package models

import (
	"database/sql"
	"time"
)

type BackupTransfer struct {
	ID                 string
	TransferDoubleHash sql.NullString
	EncryptedTransfer  string
	Recipient          string
	BlockNumber        uint64
	CreatedAt          time.Time
}

type ListOfBackupTransfer []BackupTransfer
