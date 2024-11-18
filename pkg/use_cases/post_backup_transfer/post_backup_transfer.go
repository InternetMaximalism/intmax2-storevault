package post_backup_transfer

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	postBackupTransfer "intmax2-store-vault/internal/use_cases/post_backup_transfer"

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
) postBackupTransfer.UseCasePostBackupTransfer {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context, input *postBackupTransfer.UCPostBackupTransferInput,
) (err error) {
	const (
		hName           = "UseCase PostBackupTransfer"
		transferHashKey = "transfer_hash"
		recipientKey    = "recipient"
		blockNumberKey  = "block_number"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCPostBackupTransferInputEmpty)
		return ErrUCPostBackupTransferInputEmpty
	}

	span.SetAttributes(
		attribute.String(transferHashKey, input.TransferHash),
		attribute.String(recipientKey, input.Recipient),
		attribute.Int64(blockNumberKey, int64(input.BlockNumber)),
	)

	_, err = u.db.CreateBackupTransfer(
		input.Recipient,
		input.TransferHash,
		input.EncryptedTransfer,
		int64(input.BlockNumber),
	)
	if err != nil {
		return errors.Join(ErrCreateBackupTransferWithDBFail, err)
	}

	return nil
}
