package get_backup_balances

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_balances.go -package=mocks -source=get_backup_balances.go

type UCGetBackupBalancesItem struct {
	ID                    string
	UserAddress           string
	EncryptedBalanceProof string
	EncryptedBalanceData  string
	EncryptedTxs          []string
	EncryptedTransfers    []string
	EncryptedDeposits     []string
	BlockNumber           uint64
	Signature             string
	CreatedAt             time.Time
}

type UCGetBackupBalancesMeta struct {
	StartBlockNumber uint64
	EndBlockNumber   uint64
}

type UCGetBackupBalances struct {
	Balances []*UCGetBackupBalancesItem
	Meta     *UCGetBackupBalancesMeta
}

type UCGetBackupBalancesInput struct {
	Sender           string `json:"sender"`
	StartBlockNumber uint64 `json:"startBlockNumber"`
	Limit            uint64 `json:"limit"`
}

// UseCaseGetBackupBalances describes GetBackupBalances contract.
type UseCaseGetBackupBalances interface {
	Do(ctx context.Context, input *UCGetBackupBalancesInput) (*UCGetBackupBalances, error)
}
