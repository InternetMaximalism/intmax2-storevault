package models

import (
	mFL "intmax2-store-vault/internal/sql_filter/models"
	"math/big"
	"time"
)

type BackupTransaction struct {
	ID              string
	Sender          string
	TxDoubleHash    string
	EncryptedTx     string
	EncodingVersion int64
	BlockNumber     int64
	Signature       string
	CreatedAt       time.Time
}

type ListOfBackupTransaction []BackupTransaction

type PaginationOfListOfBackupTransactionsInput struct {
	Direction mFL.Direction
	Offset    int
	Cursor    *CursorBaseOfListOfBackupTransactions
}

type CursorBaseOfListOfBackupTransactions struct {
	ID           string
	SortingValue *big.Int
}

type PaginationOfListOfBackupTransactions struct {
	Offset int
	Cursor *CursorListOfBackupTransactions
}

type CursorListOfBackupTransactions struct {
	Prev *CursorBaseOfListOfBackupTransactions
	Next *CursorBaseOfListOfBackupTransactions
}
