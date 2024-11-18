// nolint:dupl
package store_vault_server

import (
	"context"
	"fmt"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	getBackupDeposits "intmax2-store-vault/internal/use_cases/get_backup_deposits"
	"intmax2-store-vault/pkg/grpc_server/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) GetBackupDeposits(
	ctx context.Context,
	req *node.GetBackupDepositsRequest,
) (*node.GetBackupDepositsResponse, error) {
	resp := node.GetBackupDepositsResponse{}

	const (
		hName      = "Handler GetBackupDeposits"
		requestKey = "request"
		actionKey  = "action"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupDeposits.UCGetBackupDepositsInput{
		Sender:           req.Sender,
		StartBlockNumber: req.StartBlockNumber,
		Limit:            req.Limit,
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var list getBackupDeposits.UCGetBackupDeposits
	err = s.dbApp.Exec(spanCtx, &list, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var result *getBackupDeposits.UCGetBackupDeposits
		result, err = s.commands.GetBackupDeposits(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get backup deposits: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupDeposits.UCGetBackupDeposits); ok {
			v.Deposits = result.Deposits
			v.Meta = result.Meta
		} else {
			const msg = "failed to convert of the backup deposits"
			err = fmt.Errorf(msg)
			open_telemetry.MarkSpanError(spanCtx, err)
			return err
		}

		return nil
	})
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		const msg = "failed to get the backup deposits list"
		s.log.WithFields(logger.Fields{
			actionKey:  hName,
			requestKey: req.String(),
		}).WithError(err).Warnf(msg)

		return &resp, utils.OK(spanCtx)
	}

	resp.Success = true
	resp.Data = &node.GetBackupDepositsResponse_Data{
		Deposits: make([]*node.GetBackupDepositsResponse_Deposit, len(list.Deposits)),
		Meta: &node.GetBackupDepositsResponse_Meta{
			StartBlockNumber: list.Meta.StartBlockNumber,
			EndBlockNumber:   list.Meta.EndBlockNumber,
		},
	}

	for key := range list.Deposits {
		resp.Data.Deposits[key] = &node.GetBackupDepositsResponse_Deposit{
			Uuid:             list.Deposits[key].Uuid,
			Recipient:        list.Deposits[key].Recipient,
			EncryptedDeposit: list.Deposits[key].EncryptedDeposit,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: list.Deposits[key].CreatedAt.Unix(),
				Nanos:   int32(list.Deposits[key].CreatedAt.Nanosecond()),
			},
		}
	}

	return &resp, utils.OK(spanCtx)
}
