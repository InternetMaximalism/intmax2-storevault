package get_backup_transactions

import (
	"context"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_transactions.go -package=mocks -source=get_backup_transactions.go

type UCGetBackupTransactionsInput struct {
	Sender           string `json:"sender"`
	StartBlockNumber uint64 `json:"startBlockNumber"`
	Limit            uint64 `json:"limit"`
}

// UseCaseGetBackupTransactions describes GetBackupTransactions contract.
type UseCaseGetBackupTransactions interface {
	Do(ctx context.Context, input *UCGetBackupTransactionsInput) (*node.GetBackupTransactionsResponse_Data, error)
}
