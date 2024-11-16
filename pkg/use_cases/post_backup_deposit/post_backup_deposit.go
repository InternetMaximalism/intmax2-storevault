package post_backup_deposit

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	postBackupDeposit "intmax2-store-vault/internal/use_cases/post_backup_deposit"

	"go.opentelemetry.io/otel/attribute"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
}

func New(cfg *configs.Config, log logger.Logger, db SQLDriverApp) postBackupDeposit.UseCasePostBackupDeposit {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context, input *postBackupDeposit.UCPostBackupDepositInput,
) error {
	const (
		hName               = "UseCase PostBackupDeposit"
		recipientKey        = "recipient"
		blockNumberKey      = "block_number"
		depositHashKey      = "deposit_hash"
		encryptedDepositKey = "encrypted_deposit"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCPostBackupDepositInputEmpty)
		return ErrUCPostBackupDepositInputEmpty
	}

	span.SetAttributes(
		attribute.String(recipientKey, input.Recipient),
		attribute.Int64(blockNumberKey, input.BlockNumber),
		attribute.String(depositHashKey, input.DepositHash),
		attribute.String(encryptedDepositKey, input.EncryptedDeposit),
	)

	_, err := u.db.CreateBackupDeposit(
		input.Recipient, input.DepositHash, input.EncryptedDeposit, input.BlockNumber,
	)
	if err != nil {
		return errors.Join(ErrCreateBackupDepositFail, err)
	}

	return nil
}
