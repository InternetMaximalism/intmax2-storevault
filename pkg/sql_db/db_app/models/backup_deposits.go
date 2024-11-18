package models

import (
	mFL "intmax2-store-vault/internal/sql_filter/models"
	"math/big"
	"time"
)

type BackupDeposit struct {
	ID                string
	Recipient         string
	DepositDoubleHash string
	EncryptedDeposit  string
	BlockNumber       int64
	CreatedAt         time.Time
}

type ListOfBackupDeposit []BackupDeposit

type PaginationOfListOfBackupDepositsInput struct {
	Direction mFL.Direction
	Offset    int
	Cursor    *CursorBaseOfListOfBackupDeposits
}

type CursorBaseOfListOfBackupDeposits struct {
	ID           string
	SortingValue *big.Int
}

type PaginationOfListOfBackupDeposits struct {
	Offset int
	Cursor *CursorListOfBackupDeposits
}

type CursorListOfBackupDeposits struct {
	Prev *CursorBaseOfListOfBackupDeposits
	Next *CursorBaseOfListOfBackupDeposits
}
