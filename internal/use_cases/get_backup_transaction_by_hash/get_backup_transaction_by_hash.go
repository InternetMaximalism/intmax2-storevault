package get_backup_transaction_by_hash

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_transaction_by_hash.go -package=mocks -source=get_backup_transaction_by_hash.go

const (
	NotFoundMessage = "Transaction hash not found."
)

type UCGetBackupTransactionByHash struct {
	ID              string
	Sender          string
	Signature       string
	BlockNumber     uint64
	EncryptedTx     string
	EncodingVersion uint32
	CreatedAt       time.Time
}

type UCGetBackupTransactionByHashInput struct {
	Sender string `json:"sender"`
	TxHash string `json:"txHash"`
}

type UseCaseGetBackupTransactionByHash interface {
	Do(
		ctx context.Context, input *UCGetBackupTransactionByHashInput,
	) (*UCGetBackupTransactionByHash, error)
}
