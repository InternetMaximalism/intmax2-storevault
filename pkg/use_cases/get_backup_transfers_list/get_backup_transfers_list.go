package get_backup_transfers_list

import (
	"context"
	"encoding/json"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	mFL "intmax2-store-vault/internal/sql_filter/models"
	getBackupTransfersList "intmax2-store-vault/internal/use_cases/get_backup_transfers_list"
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
	"math/big"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
}

func New(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupTransfersList.UseCaseGetBackupTransfersList {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context, input *getBackupTransfersList.UCGetBackupTransfersListInput,
) (*getBackupTransfersList.UCGetBackupTransfersList, error) {
	const (
		hName    = "UseCase GetBackupTransfersList"
		inputKey = "input"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupTransfersInputEmpty)
		return nil, ErrUCGetBackupTransfersInputEmpty
	}

	bInput, err := json.Marshal(&input)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return nil, errors.Join(ErrJSONMarshalFail, err)
	}
	span.SetAttributes(attribute.String(inputKey, string(bInput)))

	var pagination mDBApp.PaginationOfListOfBackupTransfersInput
	if input.Pagination != nil {
		pagination.Direction = input.Pagination.Direction
		pagination.Offset = input.Pagination.Offset
		if input.Pagination.Cursor != nil {
			pagination.Cursor = &mDBApp.CursorBaseOfListOfBackupTransfers{
				ID:           input.Pagination.Cursor.Uuid,
				SortingValue: input.Pagination.Cursor.ConvertSortingValue,
			}
		}
	} else {
		const int100Key = 100
		pagination.Offset = int100Key
	}

	var (
		paginator *mDBApp.PaginationOfListOfBackupTransfers
		listDBApp mDBApp.ListOfBackupTransfer
	)
	paginator, listDBApp, err = u.db.GetBackupTransfersByRecipient(
		input.Recipient,
		pagination,
		input.Sorting, input.OrderBy,
		input.Filters,
	)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return nil, errors.Join(ErrGetBackupTransfersByRecipientFail, err)
	}

	resp := getBackupTransfersList.UCGetBackupTransfersList{
		List: make([]getBackupTransfersList.ItemOfGetBackupTransfersList, len(listDBApp)),
	}

	resp.Pagination = getBackupTransfersList.UCGetBackupTransfersListPaginationOfList{
		PerPage: strconv.Itoa(pagination.Offset),
	}

	if paginator.Cursor != nil {
		resp.Pagination.Cursor = &getBackupTransfersList.UCGetBackupTransfersListCursorList{}
		if paginator.Cursor.Prev != nil {
			resp.Pagination.Cursor.Prev = &getBackupTransfersList.UCGetBackupTransfersListCursorBase{
				Uuid:         paginator.Cursor.Prev.ID,
				SortingValue: paginator.Cursor.Prev.SortingValue.String(),
			}
		}
		if paginator.Cursor.Next != nil {
			resp.Pagination.Cursor.Next = &getBackupTransfersList.UCGetBackupTransfersListCursorBase{
				Uuid:         paginator.Cursor.Next.ID,
				SortingValue: paginator.Cursor.Next.SortingValue.String(),
			}
		}
	}

	for key := range listDBApp {
		resp.List[key] = getBackupTransfersList.ItemOfGetBackupTransfersList{
			Uuid:               listDBApp[key].ID,
			Recipient:          listDBApp[key].Recipient,
			TransferDoubleHash: listDBApp[key].TransferDoubleHash,
			EncryptedTransfer:  listDBApp[key].EncryptedTransfer,
			BlockNumber:        int64(listDBApp[key].BlockNumber),
			CreatedAt:          listDBApp[key].CreatedAt,
		}
		switch input.OrderBy {
		case mFL.DateCreate:
			resp.List[key].SortingValue = new(big.Int).SetInt64(listDBApp[key].CreatedAt.UTC().UnixNano()).String()
		default:
			resp.List[key].SortingValue = new(big.Int).SetInt64(listDBApp[key].CreatedAt.UTC().UnixNano()).String()
		}
	}

	return &resp, nil
}
