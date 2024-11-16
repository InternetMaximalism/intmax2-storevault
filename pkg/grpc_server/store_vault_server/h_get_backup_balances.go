package store_vault_server

import (
	"context"
	"fmt"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	getBackupBalance "intmax2-store-vault/internal/use_cases/get_backup_balances"
	"intmax2-store-vault/pkg/grpc_server/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) GetBackupBalances(
	ctx context.Context,
	req *node.GetBackupBalancesRequest,
) (*node.GetBackupBalancesResponse, error) {
	resp := node.GetBackupBalancesResponse{}

	const (
		hName      = "Handler GetBackupBalances"
		requestKey = "request"
		actionKey  = "action"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupBalance.UCGetBackupBalancesInput{
		Sender:           req.Sender,
		StartBlockNumber: req.StartBlockNumber,
		Limit:            req.Limit,
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var list getBackupBalance.UCGetBackupBalances
	err = s.dbApp.Exec(spanCtx, &list, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var result *getBackupBalance.UCGetBackupBalances
		result, err = s.commands.GetBackupBalances(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get backup balances: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupBalance.UCGetBackupBalances); ok {
			v.Balances = result.Balances
			v.Meta = result.Meta
		} else {
			const msg = "failed to convert of the backup balances"
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
	resp.Data = &node.GetBackupBalancesResponse_Data{
		Balances: make([]*node.GetBackupBalancesResponse_Balance, len(list.Balances)),
		Meta: &node.GetBackupBalancesResponse_Meta{
			StartBlockNumber: list.Meta.StartBlockNumber,
			EndBlockNumber:   list.Meta.EndBlockNumber,
		},
	}

	for key := range list.Balances {
		resp.Data.Balances[key] = &node.GetBackupBalancesResponse_Balance{
			Id:                    list.Balances[key].ID,
			UserAddress:           list.Balances[key].UserAddress,
			EncryptedBalanceProof: list.Balances[key].EncryptedBalanceProof,
			EncryptedBalanceData:  list.Balances[key].EncryptedBalanceData,
			EncryptedTxs:          list.Balances[key].EncryptedTxs,
			EncryptedTransfers:    list.Balances[key].EncryptedTransfers,
			EncryptedDeposits:     list.Balances[key].EncryptedDeposits,
			BlockNumber:           list.Balances[key].BlockNumber,
			Signature:             list.Balances[key].Signature,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: list.Balances[key].CreatedAt.Unix(),
				Nanos:   int32(list.Balances[key].CreatedAt.Nanosecond()),
			},
		}
	}

	return &resp, utils.OK(spanCtx)
}
