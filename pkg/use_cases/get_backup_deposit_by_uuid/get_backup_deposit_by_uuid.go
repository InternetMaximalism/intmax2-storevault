package get_backup_deposit_by_uuid

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	getBackupDepositByUuid "intmax2-store-vault/internal/use_cases/get_backup_deposit_by_uuid"

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
) getBackupDepositByUuid.UseCaseGetBackupDepositByUuid {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context,
	input *getBackupDepositByUuid.UCGetBackupDepositByUuidInput,
) (*getBackupDepositByUuid.UCGetBackupDepositByUuid, error) {
	const (
		hName   = "UseCase GetBackupDepositByUuid"
		uuidKey = "uuid"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupDepositByUuidInputEmpty)
		return nil, ErrUCGetBackupDepositByUuidInputEmpty
	}

	span.SetAttributes(
		attribute.String(uuidKey, input.Uuid),
	)

	deposit, err := u.db.GetBackupDepositByID(input.Uuid)
	if err != nil {
		return nil, errors.Join(ErrGetBackupDepositByIDFail, err)
	}

	return &getBackupDepositByUuid.UCGetBackupDepositByUuid{
		Uuid:             deposit.ID,
		Recipient:        deposit.Recipient,
		BlockNumber:      uint64(deposit.BlockNumber),
		EncryptedDeposit: deposit.EncryptedDeposit,
		CreatedAt:        deposit.CreatedAt,
	}, nil
}
