package store_vault_server

import (
	"context"
	"fmt"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	mFL "intmax2-store-vault/internal/sql_filter/models"
	getBackupTransfersList "intmax2-store-vault/internal/use_cases/get_backup_transfers_list"
	"intmax2-store-vault/pkg/grpc_server/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) GetBackupTransfersList(
	ctx context.Context,
	req *node.GetBackupTransfersListRequest,
) (*node.GetBackupTransfersListResponse, error) {
	resp := node.GetBackupTransfersListResponse{}

	const (
		hName      = "Handler GetBackupTransfersList"
		requestKey = "request"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupTransfersList.UCGetBackupTransfersListInput{
		Recipient: req.Recipient,
		OrderBy:   mFL.OrderBy(req.OrderBy),
		Sorting:   mFL.Sorting(req.Sorting),
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
		input.Pagination = &getBackupTransfersList.UCGetBackupTransfersListPaginationInput{
			Direction: mFL.Direction(req.Pagination.Direction),
			PerPage:   req.Pagination.PerPage,
		}

		if req.Pagination.Cursor != nil {
			input.Pagination.Cursor = &getBackupTransfersList.UCGetBackupTransfersListCursorBase{
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

	var list getBackupTransfersList.UCGetBackupTransfersList
	err = s.dbApp.Exec(spanCtx, &list, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var results *getBackupTransfersList.UCGetBackupTransfersList
		results, err = s.commands.GetBackupTransfersList(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get backup transfers: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupTransfersList.UCGetBackupTransfersList); ok {
			v.List = results.List
			v.Pagination = results.Pagination
		} else {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to convert of list of backup transfers"
			return fmt.Errorf(msg)
		}

		return nil
	})
	if err != nil {
		const msg = "failed to get backup transfers with DB App: %+v"
		return &resp, utils.Internal(spanCtx, s.log, msg, err)
	}

	resp.Success = true
	resp.Data = &node.GetBackupTransfersListResponse_Data{
		Pagination: &node.GetBackupTransfersListResponse_Pagination{
			PerPage: list.Pagination.PerPage,
		},
	}

	if list.Pagination.Cursor != nil {
		resp.Data.Pagination.Cursor = &node.GetBackupTransfersListResponse_Cursor{}
		if list.Pagination.Cursor.Prev != nil {
			resp.Data.Pagination.Cursor.Prev = &node.GetBackupTransfersListResponse_CursorBase{
				Uuid:         list.Pagination.Cursor.Prev.Uuid,
				SortingValue: list.Pagination.Cursor.Prev.SortingValue,
			}
		}
		if list.Pagination.Cursor.Next != nil {
			resp.Data.Pagination.Cursor.Next = &node.GetBackupTransfersListResponse_CursorBase{
				Uuid:         list.Pagination.Cursor.Next.Uuid,
				SortingValue: list.Pagination.Cursor.Next.SortingValue,
			}
		}
	} else if input.Pagination != nil && input.Pagination.Cursor != nil {
		resp.Data.Pagination.Cursor = &node.GetBackupTransfersListResponse_Cursor{
			Prev: &node.GetBackupTransfersListResponse_CursorBase{
				Uuid:         input.Pagination.Cursor.Uuid,
				SortingValue: input.Pagination.Cursor.SortingValue,
			},
			Next: &node.GetBackupTransfersListResponse_CursorBase{
				Uuid:         input.Pagination.Cursor.Uuid,
				SortingValue: input.Pagination.Cursor.SortingValue,
			},
		}
	}

	resp.Data.Transfers = make([]*node.GetBackupTransfersListResponse_Transfer, len(list.List))
	for key := range list.List {
		resp.Data.Transfers[key] = &node.GetBackupTransfersListResponse_Transfer{
			Uuid:              list.List[key].Uuid,
			Recipient:         list.List[key].Recipient,
			EncryptedTransfer: list.List[key].EncryptedTransfer,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: list.List[key].CreatedAt.Unix(),
				Nanos:   int32(list.List[key].CreatedAt.Nanosecond()),
			},
			SortingValue: list.List[key].SortingValue,
		}
	}

	return &resp, utils.OK(spanCtx)
}
