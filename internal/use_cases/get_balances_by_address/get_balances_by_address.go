package get_balances_by_address

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_balances_by_address.go -package=mocks -source=get_balances_by_address.go

type GetBalancesByAddressBackupDeposit struct {
	Recipient        string
	EncryptedDeposit string
	BlockNumber      uint64
	CreatedAt        time.Time
}

type GetBalancesByAddressBackupTransfer struct {
	EncryptedTransfer string
	Recipient         string
	BlockNumber       uint64
	CreatedAt         time.Time
}

type GetBalancesByAddressBackupTransaction struct {
	Sender      string
	EncryptedTx string
	BlockNumber uint64
	CreatedAt   time.Time
}

type UCGetBalancesByAddress struct {
	Deposits     []*GetBalancesByAddressBackupDeposit
	Transfers    []*GetBalancesByAddressBackupTransfer
	Transactions []*GetBalancesByAddressBackupTransaction
}

type UCGetBalancesByAddressInput struct {
	Address string `json:"address"`
}

// UseCaseGetBalancesByAddress describes GetBalancesByAddress contract.
type UseCaseGetBalancesByAddress interface {
	Do(ctx context.Context, input *UCGetBalancesByAddressInput) (*UCGetBalancesByAddress, error)
}
