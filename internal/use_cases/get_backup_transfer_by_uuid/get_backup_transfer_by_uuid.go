package get_backup_transfer_by_uuid

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_transfer_by_uuid.go -package=mocks -source=get_backup_transfer_by_uuid.go

const (
	NotFoundMessage = "Transfer not found."
)

type UCGetBackupTransferByUuid struct {
	Uuid              string
	BlockNumber       uint64
	Recipient         string
	EncryptedTransfer string
	CreatedAt         time.Time
}

type UCGetBackupTransferByUuidInput struct {
	Recipient string `json:"recipient"`
	Uuid      string `json:"uuid"`
}

type UseCaseGetBackupTransferByUuid interface {
	Do(
		ctx context.Context, input *UCGetBackupTransferByUuidInput,
	) (*UCGetBackupTransferByUuid, error)
}
