package get_backup_deposit_by_hash

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	getBackupDepositByHash "intmax2-store-vault/internal/use_cases/get_backup_deposit_by_hash"

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
) getBackupDepositByHash.UseCaseGetBackupDepositByHash {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context,
	input *getBackupDepositByHash.UCGetBackupDepositByHashInput,
) (*getBackupDepositByHash.UCGetBackupDepositByHash, error) {
	const (
		hName          = "UseCase GetBackupDepositByHash"
		senderKey      = "sender"
		depositHashKey = "deposit_hash"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupDepositByHashInputEmpty)
		return nil, ErrUCGetBackupDepositByHashInputEmpty
	}

	span.SetAttributes(
		attribute.String(depositHashKey, input.DepositHash),
		attribute.String(senderKey, input.Recipient),
	)

	deposit, err := u.db.GetBackupDepositByRecipientAndDepositDoubleHash(input.Recipient, input.DepositHash)
	if err != nil {
		return nil, errors.Join(ErrGetBackupDepositByRecipientAndDepositDoubleHash, err)
	}

	return &getBackupDepositByHash.UCGetBackupDepositByHash{
		ID:               deposit.ID,
		Recipient:        deposit.Recipient,
		BlockNumber:      uint64(deposit.BlockNumber),
		EncryptedDeposit: deposit.EncryptedDeposit,
		CreatedAt:        deposit.CreatedAt,
	}, nil
}
