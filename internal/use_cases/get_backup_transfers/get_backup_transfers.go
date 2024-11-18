package get_backup_transfers

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_transfers.go -package=mocks -source=get_backup_transfers.go

type UCGetBackupTransfersMeta struct {
	StartBlockNumber uint64
	EndBlockNumber   uint64
}

type UCGetBackupTransfersTrItem struct {
	Uuid              string
	BlockNumber       uint64
	Recipient         string
	EncryptedTransfer string
	CreatedAt         time.Time
}

type UCGetBackupTransfers struct {
	Transfers []*UCGetBackupTransfersTrItem
	Meta      *UCGetBackupTransfersMeta
}

type UCGetBackupTransfersInput struct {
	Sender           string `json:"sender"`
	StartBlockNumber uint64 `json:"startBlockNumber"`
	Limit            uint64 `json:"limit"`
}

// UseCaseGetBackupTransfers describes GetBackupTransfers contract.
type UseCaseGetBackupTransfers interface {
	Do(ctx context.Context, input *UCGetBackupTransfersInput) (*UCGetBackupTransfers, error)
}
