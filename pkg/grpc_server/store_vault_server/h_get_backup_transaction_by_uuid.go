// nolint:dupl
package store_vault_server

import (
	"context"
	"errors"
	"fmt"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	getBackupTransactionByUuid "intmax2-store-vault/internal/use_cases/get_backup_transaction_by_uuid"
	"intmax2-store-vault/pkg/grpc_server/utils"
	errorsDB "intmax2-store-vault/pkg/sql_db/errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) GetBackupTransactionByUuid(
	ctx context.Context,
	req *node.GetBackupTransactionByUuidRequest,
) (*node.GetBackupTransactionByUuidResponse, error) {
	resp := node.GetBackupTransactionByUuidResponse{}

	const (
		hName      = "Handler GetBackupTransactionByUuid"
		requestKey = "request"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupTransactionByUuid.UCGetBackupTransactionByUuidInput{
		Sender: req.Sender,
		Uuid:   req.Uuid,
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var info getBackupTransactionByUuid.UCGetBackupTransactionByUuid
	err = s.dbApp.Exec(spanCtx, &info, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var result *getBackupTransactionByUuid.UCGetBackupTransactionByUuid
		result, err = s.commands.GetBackupTransactionByUuid(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get backup transaction by uuid: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupTransactionByUuid.UCGetBackupTransactionByUuid); ok {
			v.Uuid = result.Uuid
			v.Sender = result.Sender
			v.Signature = result.Signature
			v.BlockNumber = result.BlockNumber
			v.EncryptedTx = result.EncryptedTx
			v.EncodingVersion = result.EncodingVersion
			v.CreatedAt = result.CreatedAt
		} else {
			const msg = "failed to convert of the backup transaction by uuid"
			err = fmt.Errorf(msg)
			open_telemetry.MarkSpanError(spanCtx, err)
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, errorsDB.ErrNotFound) {
			return &resp, utils.NotFound(spanCtx, fmt.Errorf("%s", getBackupTransactionByUuid.NotFoundMessage))
		}

		const msg = "failed to get backup transactions by uuid with DB App: %+v"
		return &resp, utils.Internal(spanCtx, s.log, msg, err)
	}

	resp.Success = true
	resp.Data = &node.GetBackupTransactionByUuidResponse_Data{
		Transaction: &node.GetBackupTransactionByUuidResponse_Transaction{
			Uuid:        info.Uuid,
			Sender:      info.Sender,
			Signature:   info.Signature,
			EncryptedTx: info.EncryptedTx,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: info.CreatedAt.Unix(),
				Nanos:   int32(info.CreatedAt.Nanosecond()),
			},
		},
	}

	return &resp, utils.OK(spanCtx)
}
