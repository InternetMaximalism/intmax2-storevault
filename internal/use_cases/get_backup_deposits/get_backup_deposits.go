package get_backup_deposits

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_deposits.go -package=mocks -source=get_backup_deposits.go

type UCGetBackupDepositsItem struct {
	ID               string
	Recipient        string
	BlockNumber      uint64
	EncryptedDeposit string
	CreatedAt        time.Time
}

type UCGetBackupDepositsMeta struct {
	StartBlockNumber uint64
	EndBlockNumber   uint64
}

type UCGetBackupDeposits struct {
	Deposits []*UCGetBackupDepositsItem
	Meta     *UCGetBackupDepositsMeta
}

type UCGetBackupDepositsInput struct {
	Sender           string `json:"sender"`
	StartBlockNumber uint64 `json:"startBlockNumber"`
	Limit            uint64 `json:"limit"`
}

// UseCaseGetBackupDeposits describes GetBackupDeposits contract.
type UseCaseGetBackupDeposits interface {
	Do(ctx context.Context, input *UCGetBackupDepositsInput) (*UCGetBackupDeposits, error)
}
