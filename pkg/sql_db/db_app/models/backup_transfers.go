package models

import (
	mFL "intmax2-store-vault/internal/sql_filter/models"
	"math/big"
	"time"
)

type BackupTransfer struct {
	ID                 string
	TransferDoubleHash string
	EncryptedTransfer  string
	Recipient          string
	BlockNumber        uint64
	CreatedAt          time.Time
}

type ListOfBackupTransfer []BackupTransfer

type PaginationOfListOfBackupTransfersInput struct {
	Direction mFL.Direction
	Offset    int
	Cursor    *CursorBaseOfListOfBackupTransfers
}

type CursorBaseOfListOfBackupTransfers struct {
	ID           string
	SortingValue *big.Int
}

type PaginationOfListOfBackupTransfers struct {
	Offset int
	Cursor *CursorListOfBackupTransfers
}

type CursorListOfBackupTransfers struct {
	Prev *CursorBaseOfListOfBackupTransfers
	Next *CursorBaseOfListOfBackupTransfers
}
