package store_vault_server

import (
	"context"
	"errors"
	"fmt"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	getBackupDepositByHash "intmax2-store-vault/internal/use_cases/get_backup_deposit_by_hash"
	"intmax2-store-vault/pkg/grpc_server/utils"
	errorsDB "intmax2-store-vault/pkg/sql_db/errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) GetBackupDepositByHash(
	ctx context.Context,
	req *node.GetBackupDepositByHashRequest,
) (*node.GetBackupDepositByHashResponse, error) {
	resp := node.GetBackupDepositByHashResponse{}

	const (
		hName      = "Handler GetBackupDepositByHash"
		requestKey = "request"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupDepositByHash.UCGetBackupDepositByHashInput{
		Recipient:   req.Recipient,
		DepositHash: req.DepositHash,
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var info getBackupDepositByHash.UCGetBackupDepositByHash
	err = s.dbApp.Exec(spanCtx, &info, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var result *getBackupDepositByHash.UCGetBackupDepositByHash
		result, err = s.commands.GetBackupDepositByHash(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get backup deposit by hash: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupDepositByHash.UCGetBackupDepositByHash); ok {
			v.ID = result.ID
			v.Recipient = result.Recipient
			v.BlockNumber = result.BlockNumber
			v.EncryptedDeposit = result.EncryptedDeposit
			v.CreatedAt = result.CreatedAt
		} else {
			const msg = "failed to convert of the backup deposit by hash"
			err = fmt.Errorf(msg)
			open_telemetry.MarkSpanError(spanCtx, err)
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, errorsDB.ErrNotFound) {
			return &resp, utils.NotFound(spanCtx, fmt.Errorf("%s", getBackupDepositByHash.NotFoundMessage))
		}

		const msg = "failed to get backup deposit by hash with DB App: %+v"
		return &resp, utils.Internal(spanCtx, s.log, msg, err)
	}

	resp.Success = true
	resp.Data = &node.GetBackupDepositByHashResponse_Data{
		Deposit: &node.GetBackupDepositByHashResponse_Deposit{
			Id:               info.ID,
			Recipient:        info.Recipient,
			BlockNumber:      info.BlockNumber,
			EncryptedDeposit: info.EncryptedDeposit,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: info.CreatedAt.Unix(),
				Nanos:   int32(info.CreatedAt.Nanosecond()),
			},
		},
	}

	return &resp, utils.OK(spanCtx)
}
