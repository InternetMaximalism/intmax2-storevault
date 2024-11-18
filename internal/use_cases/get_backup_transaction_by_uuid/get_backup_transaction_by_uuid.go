package get_backup_transaction_by_uuid

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_transaction_by_uuid.go -package=mocks -source=get_backup_transaction_by_uuid.go

const (
	NotFoundMessage = "Transaction not found."
)

type UCGetBackupTransactionByUuid struct {
	Uuid            string
	Sender          string
	Signature       string
	BlockNumber     uint64
	EncryptedTx     string
	EncodingVersion uint32
	CreatedAt       time.Time
}

type UCGetBackupTransactionByUuidInput struct {
	Uuid string `json:"uuid"`
}

type UseCaseGetBackupTransactionByUuid interface {
	Do(
		ctx context.Context, input *UCGetBackupTransactionByUuidInput,
	) (*UCGetBackupTransactionByUuid, error)
}
