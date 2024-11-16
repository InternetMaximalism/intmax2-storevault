package get_backup_transfer_by_hash

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	getBackupTransferByHash "intmax2-store-vault/internal/use_cases/get_backup_transfer_by_hash"

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
) getBackupTransferByHash.UseCaseGetBackupTransferByHash {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context,
	input *getBackupTransferByHash.UCGetBackupTransferByHashInput,
) (*getBackupTransferByHash.UCGetBackupTransferByHash, error) {
	const (
		hName           = "UseCase GetBackupTransferByHash"
		recipientKey    = "recipient"
		transferHashKey = "transfer_hash"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupTransferByHashInputEmpty)
		return nil, ErrUCGetBackupTransferByHashInputEmpty
	}

	span.SetAttributes(
		attribute.String(transferHashKey, input.TransferHash),
		attribute.String(recipientKey, input.Recipient),
	)

	transfer, err := u.db.GetBackupTransferByRecipientAndTransferDoubleHash(input.Recipient, input.TransferHash)
	if err != nil {
		return nil, errors.Join(ErrGetBackupTransferByRecipientAndTransferDoubleHashFail, err)
	}

	result := getBackupTransferByHash.UCGetBackupTransferByHash{
		ID:                transfer.ID,
		BlockNumber:       transfer.BlockNumber,
		Recipient:         transfer.Recipient,
		EncryptedTransfer: transfer.EncryptedTransfer,
		CreatedAt:         transfer.CreatedAt,
	}

	return &result, nil
}
