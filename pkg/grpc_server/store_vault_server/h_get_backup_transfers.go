// nolint:dupl
package store_vault_server

import (
	"context"
	"fmt"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	getBackupTransfers "intmax2-store-vault/internal/use_cases/get_backup_transfers"
	"intmax2-store-vault/pkg/grpc_server/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) GetBackupTransfers(
	ctx context.Context,
	req *node.GetBackupTransfersRequest,
) (*node.GetBackupTransfersResponse, error) {
	resp := node.GetBackupTransfersResponse{}

	const (
		hName      = "Handler GetBackupTransfers"
		requestKey = "request"
		actionKey  = "action"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupTransfers.UCGetBackupTransfersInput{
		Sender:           req.Sender,
		StartBlockNumber: req.StartBlockNumber,
		Limit:            req.Limit,
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var list getBackupTransfers.UCGetBackupTransfers
	err = s.dbApp.Exec(spanCtx, &list, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var result *getBackupTransfers.UCGetBackupTransfers
		result, err = s.commands.GetBackupTransfers(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get backup transfers: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupTransfers.UCGetBackupTransfers); ok {
			v.Transfers = result.Transfers
			v.Meta = result.Meta
		} else {
			const msg = "failed to convert of the backup transfers"
			err = fmt.Errorf(msg)
			open_telemetry.MarkSpanError(spanCtx, err)
			return err
		}

		return nil
	})
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		const msg = "failed to get the backup transfers list"
		s.log.WithFields(logger.Fields{
			actionKey:  hName,
			requestKey: req.String(),
		}).WithError(err).Warnf(msg)

		return &resp, utils.OK(spanCtx)
	}

	resp.Success = true
	resp.Data = &node.GetBackupTransfersResponse_Data{
		Transfers: make([]*node.GetBackupTransfersResponse_Transfer, len(list.Transfers)),
		Meta: &node.GetBackupTransfersResponse_Meta{
			StartBlockNumber: list.Meta.StartBlockNumber,
			EndBlockNumber:   list.Meta.EndBlockNumber,
		},
	}

	for key := range list.Transfers {
		resp.Data.Transfers[key] = &node.GetBackupTransfersResponse_Transfer{
			Uuid:              list.Transfers[key].Uuid,
			Recipient:         list.Transfers[key].Recipient,
			EncryptedTransfer: list.Transfers[key].EncryptedTransfer,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: list.Transfers[key].CreatedAt.Unix(),
				Nanos:   int32(list.Transfers[key].CreatedAt.Nanosecond()),
			},
		}
	}

	return &resp, utils.OK(spanCtx)
}
