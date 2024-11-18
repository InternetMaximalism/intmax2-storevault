// nolint:dupl
package store_vault_server

import (
	"context"
	"errors"
	"fmt"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	getBackupTransferByUuid "intmax2-store-vault/internal/use_cases/get_backup_transfer_by_uuid"
	"intmax2-store-vault/pkg/grpc_server/utils"
	errorsDB "intmax2-store-vault/pkg/sql_db/errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) GetBackupTransferByUuid(
	ctx context.Context,
	req *node.GetBackupTransferByUuidRequest,
) (*node.GetBackupTransferByUuidResponse, error) {
	resp := node.GetBackupTransferByUuidResponse{}

	const (
		hName      = "Handler GetBackupTransferByUuid"
		requestKey = "request"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupTransferByUuid.UCGetBackupTransferByUuidInput{
		Uuid: req.Uuid,
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var info getBackupTransferByUuid.UCGetBackupTransferByUuid
	err = s.dbApp.Exec(spanCtx, &info, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var result *getBackupTransferByUuid.UCGetBackupTransferByUuid
		result, err = s.commands.GetBackupTransferByUuid(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get backup transfer by uuid: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupTransferByUuid.UCGetBackupTransferByUuid); ok {
			v.Uuid = result.Uuid
			v.BlockNumber = result.BlockNumber
			v.Recipient = result.Recipient
			v.EncryptedTransfer = result.EncryptedTransfer
			v.CreatedAt = result.CreatedAt
		} else {
			const msg = "failed to convert of the backup transfer by uuid"
			err = fmt.Errorf(msg)
			open_telemetry.MarkSpanError(spanCtx, err)
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, errorsDB.ErrNotFound) {
			return &resp, utils.NotFound(spanCtx, fmt.Errorf("%s", getBackupTransferByUuid.NotFoundMessage))
		}

		const msg = "failed to get backup transfer by uuid with DB App: %+v"
		return &resp, utils.Internal(spanCtx, s.log, msg, err)
	}

	resp.Success = true
	resp.Data = &node.GetBackupTransferByUuidResponse_Data{
		Transfer: &node.GetBackupTransferByUuidResponse_Transfer{
			Uuid:              info.Uuid,
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
