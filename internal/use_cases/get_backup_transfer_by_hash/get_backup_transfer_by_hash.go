package get_backup_transfer_by_hash

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_transfer_by_hash.go -package=mocks -source=get_backup_transfer_by_hash.go

const (
	NotFoundMessage = "Transfer hash not found."
)

type UCGetBackupTransferByHash struct {
	ID                string
	BlockNumber       uint64
	Recipient         string
	EncryptedTransfer string
	CreatedAt         time.Time
}

type UCGetBackupTransferByHashInput struct {
	Recipient    string `json:"recipient"`
	TransferHash string `json:"transferHash"`
}

type UseCaseGetBackupTransferByHash interface {
	Do(
		ctx context.Context, input *UCGetBackupTransferByHashInput,
	) (*UCGetBackupTransferByHash, error)
}
