// nolint:dupl
package store_vault_server

import (
	"context"
	"fmt"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	getBackupTransactions "intmax2-store-vault/internal/use_cases/get_backup_transactions"
	"intmax2-store-vault/pkg/grpc_server/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) GetBackupTransactions(
	ctx context.Context,
	req *node.GetBackupTransactionsRequest,
) (*node.GetBackupTransactionsResponse, error) {
	resp := node.GetBackupTransactionsResponse{}

	const (
		hName      = "Handler GetBackupTransactions"
		requestKey = "request"
		actionKey  = "action"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupTransactions.UCGetBackupTransactionsInput{
		Sender:           req.Sender,
		StartBlockNumber: req.StartBlockNumber,
		Limit:            req.Limit,
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var list getBackupTransactions.UCGetBackupTransactions
	err = s.dbApp.Exec(spanCtx, &list, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var result *getBackupTransactions.UCGetBackupTransactions
		result, err = s.commands.GetBackupTransactions(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get backup transactions: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupTransactions.UCGetBackupTransactions); ok {
			v.Transactions = result.Transactions
			v.Meta = result.Meta
		} else {
			const msg = "failed to convert of the backup transactions"
			err = fmt.Errorf(msg)
			open_telemetry.MarkSpanError(spanCtx, err)
			return err
		}

		return nil
	})
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		const msg = "failed to get the backup transactions list"
		s.log.WithFields(logger.Fields{
			actionKey:  hName,
			requestKey: req.String(),
		}).WithError(err).Warnf(msg)

		return &resp, utils.OK(spanCtx)
	}

	resp.Success = true
	resp.Data = &node.GetBackupTransactionsResponse_Data{
		Transactions: make([]*node.GetBackupTransactionsResponse_Transaction, len(list.Transactions)),
		Meta: &node.GetBackupTransactionsResponse_Meta{
			StartBlockNumber: list.Meta.StartBlockNumber,
			EndBlockNumber:   list.Meta.EndBlockNumber,
		},
	}

	for key := range list.Transactions {
		resp.Data.Transactions[key] = &node.GetBackupTransactionsResponse_Transaction{
			Uuid:        list.Transactions[key].Uuid,
			Sender:      list.Transactions[key].Sender,
			Signature:   list.Transactions[key].Signature,
			EncryptedTx: list.Transactions[key].EncryptedTx,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: list.Transactions[key].CreatedAt.Unix(),
				Nanos:   int32(list.Transactions[key].CreatedAt.Nanosecond()),
			},
		}
	}

	return &resp, utils.OK(spanCtx)
}
