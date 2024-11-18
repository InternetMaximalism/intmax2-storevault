// nolint:dupl
package store_vault_server

import (
	"context"
	"fmt"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	mFL "intmax2-store-vault/internal/sql_filter/models"
	getBackupTransactionsList "intmax2-store-vault/internal/use_cases/get_backup_transactions_list"
	"intmax2-store-vault/pkg/grpc_server/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) GetBackupTransactionsList(
	ctx context.Context,
	req *node.GetBackupTransactionsListRequest,
) (*node.GetBackupTransactionsListResponse, error) {
	resp := node.GetBackupTransactionsListResponse{}

	const (
		hName      = "Handler GetBackupTransactionsList"
		requestKey = "request"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupTransactionsList.UCGetBackupTransactionsListInput{
		Sender:  req.Sender,
		OrderBy: mFL.OrderBy(req.OrderBy),
		Sorting: mFL.Sorting(req.Sorting),
	}

	input.Filters = make([]*mFL.Filter, len(req.Filter))
	for i := range req.Filter {
		input.Filters[i] = &mFL.Filter{
			Relation:  mFL.Relation(req.Filter[i].Relation),
			DataField: mFL.DataField(req.Filter[i].DataField),
			Condition: mFL.Condition(req.Filter[i].Condition),
			Value:     req.Filter[i].Value,
		}
	}

	if req.Pagination != nil {
		input.Pagination = &getBackupTransactionsList.UCGetBackupTransactionsListPaginationInput{
			Direction: mFL.Direction(req.Pagination.Direction),
			PerPage:   req.Pagination.PerPage,
		}

		if req.Pagination.Cursor != nil {
			input.Pagination.Cursor = &getBackupTransactionsList.UCGetBackupTransactionsListCursorBase{
				Uuid:         req.Pagination.Cursor.Uuid,
				SortingValue: req.Pagination.Cursor.SortingValue,
			}
		}
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var list getBackupTransactionsList.UCGetBackupTransactionsList
	err = s.dbApp.Exec(spanCtx, &list, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var results *getBackupTransactionsList.UCGetBackupTransactionsList
		results, err = s.commands.GetBackupTransactionsList(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get backup transactions: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupTransactionsList.UCGetBackupTransactionsList); ok {
			v.List = results.List
			v.Pagination = results.Pagination
		} else {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to convert of list of backup transactions"
			return fmt.Errorf(msg)
		}

		return nil
	})
	if err != nil {
		const msg = "failed to get backup transactions with DB App: %+v"
		return &resp, utils.Internal(spanCtx, s.log, msg, err)
	}

	resp.Success = true
	resp.Data = &node.GetBackupTransactionsListResponse_Data{
		Pagination: &node.GetBackupTransactionsListResponse_Pagination{
			PerPage: list.Pagination.PerPage,
		},
	}

	if list.Pagination.Cursor != nil {
		resp.Data.Pagination.Cursor = &node.GetBackupTransactionsListResponse_Cursor{}
		if list.Pagination.Cursor.Prev != nil {
			resp.Data.Pagination.Cursor.Prev = &node.GetBackupTransactionsListResponse_CursorBase{
				Uuid:         list.Pagination.Cursor.Prev.Uuid,
				SortingValue: list.Pagination.Cursor.Prev.SortingValue,
			}
		}
		if list.Pagination.Cursor.Next != nil {
			resp.Data.Pagination.Cursor.Next = &node.GetBackupTransactionsListResponse_CursorBase{
				Uuid:         list.Pagination.Cursor.Next.Uuid,
				SortingValue: list.Pagination.Cursor.Next.SortingValue,
			}
		}
	} else if input.Pagination != nil && input.Pagination.Cursor != nil {
		resp.Data.Pagination.Cursor = &node.GetBackupTransactionsListResponse_Cursor{
			Prev: &node.GetBackupTransactionsListResponse_CursorBase{
				Uuid:         input.Pagination.Cursor.Uuid,
				SortingValue: input.Pagination.Cursor.SortingValue,
			},
			Next: &node.GetBackupTransactionsListResponse_CursorBase{
				Uuid:         input.Pagination.Cursor.Uuid,
				SortingValue: input.Pagination.Cursor.SortingValue,
			},
		}
	}

	resp.Data.Transactions = make([]*node.GetBackupTransactionsListResponse_Transaction, len(list.List))
	for key := range list.List {
		resp.Data.Transactions[key] = &node.GetBackupTransactionsListResponse_Transaction{
			Uuid:        list.List[key].Uuid,
			Sender:      list.List[key].Sender,
			Signature:   list.List[key].Signature,
			EncryptedTx: list.List[key].EncryptedTx,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: list.List[key].CreatedAt.Unix(),
				Nanos:   int32(list.List[key].CreatedAt.Nanosecond()),
			},
			SortingValue: list.List[key].SortingValue,
		}
	}

	return &resp, utils.OK(spanCtx)
}
