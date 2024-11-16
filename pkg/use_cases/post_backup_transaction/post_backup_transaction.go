package post_backup_transaction

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	postBackupTransaction "intmax2-store-vault/internal/use_cases/post_backup_transaction"

	"go.opentelemetry.io/otel/attribute"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
}

func New(cfg *configs.Config, log logger.Logger, db SQLDriverApp) postBackupTransaction.UseCasePostBackupTransaction {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context, input *postBackupTransaction.UCPostBackupTransactionInput,
) error {
	const (
		hName          = "UseCase PostBackupTransaction"
		senderKey      = "sender"
		blockNumberKey = "block_number"
		txHashKey      = "tx_hash"
		encryptedTxKey = "encrypted_tx"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCPostBackupTransactionInputEmpty)
		return ErrUCPostBackupTransactionInputEmpty
	}

	span.SetAttributes(
		attribute.String(senderKey, input.Sender),
		attribute.Int64(blockNumberKey, int64(input.BlockNumber)),
		attribute.String(txHashKey, input.TxHash),
		attribute.String(encryptedTxKey, input.EncryptedTx),
	)

	_, err := u.db.CreateBackupTransaction(
		input.Sender, input.TxHash, input.EncryptedTx, input.Signature, int64(input.BlockNumber),
	)
	if err != nil {
		return errors.Join(ErrCreateBackupTransactionFail, err)
	}

	return nil
}
