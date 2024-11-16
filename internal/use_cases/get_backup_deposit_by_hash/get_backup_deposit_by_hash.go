package get_backup_deposit_by_hash

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_deposit_by_hash.go -package=mocks -source=get_backup_deposit_by_hash.go

const (
	NotFoundMessage = "Deposit hash not found."
)

type UCGetBackupDepositByHash struct {
	ID               string
	Recipient        string
	BlockNumber      uint64
	EncryptedDeposit string
	CreatedAt        time.Time
}

type UCGetBackupDepositByHashInput struct {
	Recipient   string `json:"recipient"`
	DepositHash string `json:"depositHash"`
}

type UseCaseGetBackupDepositByHash interface {
	Do(
		ctx context.Context, input *UCGetBackupDepositByHashInput,
	) (*UCGetBackupDepositByHash, error)
}
