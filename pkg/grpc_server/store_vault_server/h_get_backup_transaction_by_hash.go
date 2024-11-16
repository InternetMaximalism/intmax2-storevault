package store_vault_server

import (
	"context"
	"errors"
	"fmt"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	getBackupTransactionByHash "intmax2-store-vault/internal/use_cases/get_backup_transaction_by_hash"
	"intmax2-store-vault/pkg/grpc_server/utils"
	errorsDB "intmax2-store-vault/pkg/sql_db/errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) GetBackupTransactionByHash(
	ctx context.Context,
	req *node.GetBackupTransactionByHashRequest,
) (*node.GetBackupTransactionByHashResponse, error) {
	resp := node.GetBackupTransactionByHashResponse{}

	const (
		hName      = "Handler GetBackupTransactionByHash"
		requestKey = "request"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupTransactionByHash.UCGetBackupTransactionByHashInput{
		Sender: req.Sender,
		TxHash: req.TxHash,
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var info getBackupTransactionByHash.UCGetBackupTransactionByHash
	err = s.dbApp.Exec(spanCtx, &info, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var result *getBackupTransactionByHash.UCGetBackupTransactionByHash
		result, err = s.commands.GetBackupTransactionByHash(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get backup transaction by hash: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupTransactionByHash.UCGetBackupTransactionByHash); ok {
			v.ID = result.ID
			v.Sender = result.Sender
			v.Signature = result.Signature
			v.BlockNumber = result.BlockNumber
			v.EncryptedTx = result.EncryptedTx
			v.EncodingVersion = result.EncodingVersion
			v.CreatedAt = result.CreatedAt
		} else {
			const msg = "failed to convert of the backup transaction by hash"
			err = fmt.Errorf(msg)
			open_telemetry.MarkSpanError(spanCtx, err)
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, errorsDB.ErrNotFound) {
			return &resp, utils.NotFound(spanCtx, fmt.Errorf("%s", getBackupTransactionByHash.NotFoundMessage))
		}

		const msg = "failed to get backup transactions with DB App: %+v"
		return &resp, utils.Internal(spanCtx, s.log, msg, err)
	}

	resp.Success = true
	resp.Data = &node.GetBackupTransactionByHashResponse_Data{
		Transaction: &node.GetBackupTransactionByHashResponse_Transaction{
			Id:              info.ID,
			Sender:          info.Sender,
			Signature:       info.Signature,
			BlockNumber:     info.BlockNumber,
			EncryptedTx:     info.EncryptedTx,
			EncodingVersion: info.EncodingVersion,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: info.CreatedAt.Unix(),
				Nanos:   int32(info.CreatedAt.Nanosecond()),
			},
		},
	}

	return &resp, utils.OK(spanCtx)
}
