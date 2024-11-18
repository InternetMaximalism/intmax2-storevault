package store_vault_server

import (
	"context"
	"fmt"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	getBackupBalancesByAddress "intmax2-store-vault/internal/use_cases/get_balances_by_address"
	"intmax2-store-vault/pkg/grpc_server/utils"
	"math/big"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) GetBalancesByAddress(
	ctx context.Context,
	req *node.GetBalancesByAddressRequest,
) (*node.GetBalancesByAddressResponse, error) {
	resp := node.GetBalancesByAddressResponse{}

	const (
		hName      = "Handler GetBalancesByAddress"
		requestKey = "request"
		actionKey  = "action"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupBalancesByAddress.UCGetBalancesByAddressInput{
		Address: req.Address,
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var list getBackupBalancesByAddress.UCGetBalancesByAddress
	err = s.dbApp.Exec(spanCtx, &list, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var result *getBackupBalancesByAddress.UCGetBalancesByAddress
		result, err = s.commands.GetBalancesByAddress(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get balances by address request: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupBalancesByAddress.UCGetBalancesByAddress); ok {
			v.Deposits = result.Deposits
			v.Transfers = result.Transfers
			v.Transactions = result.Transactions
		} else {
			const msg = "failed to convert of the balances by address"
			err = fmt.Errorf(msg)
			open_telemetry.MarkSpanError(spanCtx, err)
			return err
		}

		return nil
	})
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		const msg = "failed to get the balances by address"
		s.log.WithFields(logger.Fields{
			actionKey:  hName,
			requestKey: req.String(),
		}).WithError(err).Warnf(msg)

		return &resp, utils.OK(spanCtx)
	}

	resp.Success = true
	resp.Deposits = make([]*node.BackupDeposit, len(list.Deposits))
	resp.Transfers = make([]*node.BackupTransfer, len(list.Transfers))
	resp.Transactions = make([]*node.BackupTransaction, len(list.Transactions))

	for key := range list.Deposits {
		resp.Deposits[key] = &node.BackupDeposit{
			Recipient:        list.Deposits[key].Recipient,
			EncryptedDeposit: list.Deposits[key].EncryptedDeposit,
			BlockNumber:      list.Deposits[key].BlockNumber,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: list.Deposits[key].CreatedAt.Unix(),
				Nanos:   int32(list.Deposits[key].CreatedAt.Nanosecond()),
			},
		}
	}

	for key := range list.Transfers {
		resp.Transfers[key] = &node.BackupTransfer{
			EncryptedTransfer: list.Transfers[key].EncryptedTransfer,
			Recipient:         list.Transfers[key].Recipient,
			BlockNumber:       list.Transfers[key].BlockNumber,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: list.Transfers[key].CreatedAt.Unix(),
				Nanos:   int32(list.Transfers[key].CreatedAt.Nanosecond()),
			},
		}
	}

	for key := range list.Transactions {
		resp.Transactions[key] = &node.BackupTransaction{
			Sender:      list.Transactions[key].Sender,
			EncryptedTx: list.Transactions[key].EncryptedTx,
			BlockNumber: new(big.Int).SetUint64(list.Transactions[key].BlockNumber).String(),
			CreatedAt: &timestamppb.Timestamp{
				Seconds: list.Transactions[key].CreatedAt.Unix(),
				Nanos:   int32(list.Transactions[key].CreatedAt.Nanosecond()),
			},
		}
	}

	return &resp, utils.OK(spanCtx)
}
