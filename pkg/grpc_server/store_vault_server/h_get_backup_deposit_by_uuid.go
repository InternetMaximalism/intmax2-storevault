// nolint:dupl
package store_vault_server

import (
	"context"
	"errors"
	"fmt"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	getBackupDepositByUuid "intmax2-store-vault/internal/use_cases/get_backup_deposit_by_uuid"
	"intmax2-store-vault/pkg/grpc_server/utils"
	errorsDB "intmax2-store-vault/pkg/sql_db/errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) GetBackupDepositByUuid(
	ctx context.Context,
	req *node.GetBackupDepositByUuidRequest,
) (*node.GetBackupDepositByUuidResponse, error) {
	resp := node.GetBackupDepositByUuidResponse{}

	const (
		hName      = "Handler GetBackupDepositByUuid"
		requestKey = "request"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupDepositByUuid.UCGetBackupDepositByUuidInput{
		Uuid: req.Uuid,
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var info getBackupDepositByUuid.UCGetBackupDepositByUuid
	err = s.dbApp.Exec(spanCtx, &info, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var result *getBackupDepositByUuid.UCGetBackupDepositByUuid
		result, err = s.commands.GetBackupDepositByUuid(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get backup deposit by uuid: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupDepositByUuid.UCGetBackupDepositByUuid); ok {
			v.Uuid = result.Uuid
			v.Recipient = result.Recipient
			v.BlockNumber = result.BlockNumber
			v.EncryptedDeposit = result.EncryptedDeposit
			v.CreatedAt = result.CreatedAt
		} else {
			const msg = "failed to convert of the backup deposit by uuid"
			err = fmt.Errorf(msg)
			open_telemetry.MarkSpanError(spanCtx, err)
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, errorsDB.ErrNotFound) {
			return &resp, utils.NotFound(spanCtx, fmt.Errorf("%s", getBackupDepositByUuid.NotFoundMessage))
		}

		const msg = "failed to get backup deposit by uuid with DB App: %+v"
		return &resp, utils.Internal(spanCtx, s.log, msg, err)
	}

	resp.Success = true
	resp.Data = &node.GetBackupDepositByUuidResponse_Data{
		Deposit: &node.GetBackupDepositByUuidResponse_Deposit{
			Uuid:             info.Uuid,
			Recipient:        info.Recipient,
			EncryptedDeposit: info.EncryptedDeposit,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: info.CreatedAt.Unix(),
				Nanos:   int32(info.CreatedAt.Nanosecond()),
			},
		},
	}

	return &resp, utils.OK(spanCtx)
}
