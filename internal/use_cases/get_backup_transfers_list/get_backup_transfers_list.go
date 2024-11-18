package get_backup_transfers_list

import (
	"context"
	mFL "intmax2-store-vault/internal/sql_filter/models"
	"math/big"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_transfers_list.go -package=mocks -source=get_backup_transfers_list.go

type UCGetBackupTransfersListCursorBase struct {
	Uuid                string   `json:"uuid"`
	SortingValue        string   `json:"sorting_value"`
	ConvertSortingValue *big.Int `json:"-"`
}

type UCGetBackupTransfersListPaginationInput struct {
	Direction mFL.Direction                       `json:"direction"`
	PerPage   string                              `json:"per_page"`
	Offset    int                                 `json:"-"`
	Cursor    *UCGetBackupTransfersListCursorBase `json:"cursor"`
}

type UCGetBackupTransfersListInput struct {
	Recipient  string                                   `json:"recipient"`
	Pagination *UCGetBackupTransfersListPaginationInput `json:"pagination"`
	OrderBy    mFL.OrderBy                              `json:"order_by"`
	Sorting    mFL.Sorting                              `json:"sorting"`
	Filters    []*mFL.Filter                            `json:"filters"`
}

type UCGetBackupTransfersListCursorList struct {
	Prev *UCGetBackupTransfersListCursorBase
	Next *UCGetBackupTransfersListCursorBase
}

type UCGetBackupTransfersListPaginationOfList struct {
	PerPage string
	Cursor  *UCGetBackupTransfersListCursorList
}

type ItemOfGetBackupTransfersList struct {
	Uuid               string
	Recipient          string
	TransferDoubleHash string
	EncryptedTransfer  string
	BlockNumber        int64
	CreatedAt          time.Time
	SortingValue       string
}

type UCGetBackupTransfersList struct {
	Pagination UCGetBackupTransfersListPaginationOfList
	List       []ItemOfGetBackupTransfersList
}

// UseCaseGetBackupTransfersList describes GetBackupTransfersList contract.
type UseCaseGetBackupTransfersList interface {
	Do(ctx context.Context, input *UCGetBackupTransfersListInput) (*UCGetBackupTransfersList, error)
}
