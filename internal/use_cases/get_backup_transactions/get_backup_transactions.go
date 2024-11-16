package get_backup_transactions

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_transactions.go -package=mocks -source=get_backup_transactions.go

type UCGeyBackupTransactionsItem struct {
	ID              string
	Sender          string
	Signature       string
	BlockNumber     uint64
	EncryptedTx     string
	EncodingVersion uint32
	CreatedAt       time.Time
}

type UCGetBackupTransactionsMeta struct {
	StartBlockNumber uint64
	EndBlockNumber   uint64
}

type UCGetBackupTransactions struct {
	Transactions []*UCGeyBackupTransactionsItem
	Meta         *UCGetBackupTransactionsMeta
}

type UCGetBackupTransactionsInput struct {
	Sender           string `json:"sender"`
	StartBlockNumber uint64 `json:"startBlockNumber"`
	Limit            uint64 `json:"limit"`
}

// UseCaseGetBackupTransactions describes GetBackupTransactions contract.
type UseCaseGetBackupTransactions interface {
	Do(ctx context.Context, input *UCGetBackupTransactionsInput) (*UCGetBackupTransactions, error)
}
