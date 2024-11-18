package get_backup_deposits_list

import (
	"context"
	mFL "intmax2-store-vault/internal/sql_filter/models"
	"math/big"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_deposits_list.go -package=mocks -source=get_backup_deposits_list.go

type UCGetBackupDepositsListCursorBase struct {
	Uuid                string   `json:"uuid"`
	SortingValue        string   `json:"sorting_value"`
	ConvertSortingValue *big.Int `json:"-"`
}

type UCGetBackupDepositsListPaginationInput struct {
	Direction mFL.Direction                      `json:"direction"`
	PerPage   string                             `json:"per_page"`
	Offset    int                                `json:"-"`
	Cursor    *UCGetBackupDepositsListCursorBase `json:"cursor"`
}

type UCGetBackupDepositsListInput struct {
	Recipient  string                                  `json:"recipient"`
	Pagination *UCGetBackupDepositsListPaginationInput `json:"pagination"`
	OrderBy    mFL.OrderBy                             `json:"order_by"`
	Sorting    mFL.Sorting                             `json:"sorting"`
	Filters    []*mFL.Filter                           `json:"filters"`
}

type UCGetBackupDepositsListCursorList struct {
	Prev *UCGetBackupDepositsListCursorBase
	Next *UCGetBackupDepositsListCursorBase
}

type UCGetBackupDepositsListPaginationOfList struct {
	PerPage string
	Cursor  *UCGetBackupDepositsListCursorList
}

type ItemOfGetBackupDepositsList struct {
	Uuid              string
	Recipient         string
	DepositDoubleHash string
	EncryptedDeposit  string
	BlockNumber       int64
	CreatedAt         time.Time
	SortingValue      string
}

type UCGetBackupDepositsList struct {
	Pagination UCGetBackupDepositsListPaginationOfList
	List       []ItemOfGetBackupDepositsList
}

// UseCaseGetBackupDepositsList describes GetBackupDepositsList contract.
type UseCaseGetBackupDepositsList interface {
	Do(ctx context.Context, input *UCGetBackupDepositsListInput) (*UCGetBackupDepositsList, error)
}
