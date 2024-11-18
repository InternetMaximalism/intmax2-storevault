package get_backup_transfer_by_uuid

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	getBackupTransferByUuid "intmax2-store-vault/internal/use_cases/get_backup_transfer_by_uuid"

	"go.opentelemetry.io/otel/attribute"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
}

func New(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupTransferByUuid.UseCaseGetBackupTransferByUuid {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context,
	input *getBackupTransferByUuid.UCGetBackupTransferByUuidInput,
) (*getBackupTransferByUuid.UCGetBackupTransferByUuid, error) {
	const (
		hName   = "UseCase GetBackupTransferByUuid"
		uuidKey = "uuid"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupTransferByUuidInputEmpty)
		return nil, ErrUCGetBackupTransferByUuidInputEmpty
	}

	span.SetAttributes(
		attribute.String(uuidKey, input.Uuid),
	)

	transfer, err := u.db.GetBackupTransferByID(input.Uuid)
	if err != nil {
		return nil, errors.Join(ErrGetBackupTransferByIDFail, err)
	}

	result := getBackupTransferByUuid.UCGetBackupTransferByUuid{
		Uuid:              transfer.ID,
		BlockNumber:       transfer.BlockNumber,
		Recipient:         transfer.Recipient,
		EncryptedTransfer: transfer.EncryptedTransfer,
		CreatedAt:         transfer.CreatedAt,
	}

	return &result, nil
}
