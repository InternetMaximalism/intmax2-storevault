package get_backup_transactions_list

import (
	"context"
	mFL "intmax2-store-vault/internal/sql_filter/models"
	"math/big"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_transactions_list.go -package=mocks -source=get_backup_transactions_list.go

type UCGetBackupTransactionsListCursorBase struct {
	Uuid                string   `json:"uuid"`
	SortingValue        string   `json:"sorting_value"`
	ConvertSortingValue *big.Int `json:"-"`
}

type UCGetBackupTransactionsListPaginationInput struct {
	Direction mFL.Direction                          `json:"direction"`
	PerPage   string                                 `json:"per_page"`
	Offset    int                                    `json:"-"`
	Cursor    *UCGetBackupTransactionsListCursorBase `json:"cursor"`
}

type UCGetBackupTransactionsListInput struct {
	Sender     string                                      `json:"sender"`
	Pagination *UCGetBackupTransactionsListPaginationInput `json:"pagination"`
	OrderBy    mFL.OrderBy                                 `json:"order_by"`
	Sorting    mFL.Sorting                                 `json:"sorting"`
	Filters    []*mFL.Filter                               `json:"filters"`
}

type UCGetBackupTransactionsListCursorList struct {
	Prev *UCGetBackupTransactionsListCursorBase
	Next *UCGetBackupTransactionsListCursorBase
}

type UCGetBackupTransactionsListPaginationOfList struct {
	PerPage string
	Cursor  *UCGetBackupTransactionsListCursorList
}

type ItemOfGetBackupTransactionsList struct {
	Uuid         string
	Sender       string
	EncryptedTx  string
	Signature    string
	CreatedAt    time.Time
	SortingValue string
}

type UCGetBackupTransactionsList struct {
	Pagination UCGetBackupTransactionsListPaginationOfList
	List       []ItemOfGetBackupTransactionsList
}

// UseCaseGetBackupTransactionsList describes GetBackupTransactionsList contract.
type UseCaseGetBackupTransactionsList interface {
	Do(ctx context.Context, input *UCGetBackupTransactionsListInput) (*UCGetBackupTransactionsList, error)
}
