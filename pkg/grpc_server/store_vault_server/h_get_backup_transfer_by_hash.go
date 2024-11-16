package store_vault_server

import (
	"context"
	"errors"
	"fmt"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	getBackupTransferByHash "intmax2-store-vault/internal/use_cases/get_backup_transfer_by_hash"
	"intmax2-store-vault/pkg/grpc_server/utils"
	errorsDB "intmax2-store-vault/pkg/sql_db/errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) GetBackupTransferByHash(
	ctx context.Context,
	req *node.GetBackupTransferByHashRequest,
) (*node.GetBackupTransferByHashResponse, error) {
	resp := node.GetBackupTransferByHashResponse{}

	const (
		hName      = "Handler GetBackupTransferByHash"
		requestKey = "request"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupTransferByHash.UCGetBackupTransferByHashInput{
		Recipient:    req.Recipient,
		TransferHash: req.TransferHash,
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var info getBackupTransferByHash.UCGetBackupTransferByHash
	err = s.dbApp.Exec(spanCtx, &info, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var result *getBackupTransferByHash.UCGetBackupTransferByHash
		result, err = s.commands.GetBackupTransferByHash(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get backup transfer by hash: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupTransferByHash.UCGetBackupTransferByHash); ok {
			v.ID = result.ID
			v.BlockNumber = result.BlockNumber
			v.Recipient = result.Recipient
			v.EncryptedTransfer = result.EncryptedTransfer
			v.CreatedAt = result.CreatedAt
		} else {
			const msg = "failed to convert of the backup transfer by hash"
			err = fmt.Errorf(msg)
			open_telemetry.MarkSpanError(spanCtx, err)
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, errorsDB.ErrNotFound) {
			return &resp, utils.NotFound(spanCtx, fmt.Errorf("%s", getBackupTransferByHash.NotFoundMessage))
		}

		const msg = "failed to get backup transfer with DB App: %+v"
		return &resp, utils.Internal(spanCtx, s.log, msg, err)
	}

	resp.Success = true
	resp.Data = &node.GetBackupTransferByHashResponse_Data{
		Transfer: &node.GetBackupTransferByHashResponse_Transfer{
			Id:                info.ID,
			BlockNumber:       info.BlockNumber,
			Recipient:         info.Recipient,
			EncryptedTransfer: info.EncryptedTransfer,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: info.CreatedAt.Unix(),
				Nanos:   int32(info.CreatedAt.Nanosecond()),
			},
		},
	}

	return &resp, utils.OK(spanCtx)
}
