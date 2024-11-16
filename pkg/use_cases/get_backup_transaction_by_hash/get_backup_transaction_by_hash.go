package get_backup_transaction_by_hash

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	getBackupTransactionByHash "intmax2-store-vault/internal/use_cases/get_backup_transaction_by_hash"

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
) getBackupTransactionByHash.UseCaseGetBackupTransactionByHash {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context,
	input *getBackupTransactionByHash.UCGetBackupTransactionByHashInput,
) (*getBackupTransactionByHash.UCGetBackupTransactionByHash, error) {
	const (
		hName     = "UseCase GetBackupTransactionByHash"
		senderKey = "sender"
		txHashKey = "tx_hash"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupTransactionByHashInputEmpty)
		return nil, ErrUCGetBackupTransactionByHashInputEmpty
	}

	span.SetAttributes(
		attribute.String(txHashKey, input.TxHash),
		attribute.String(senderKey, input.Sender),
	)

	transaction, err := u.db.GetBackupTransactionBySenderAndTxDoubleHash(input.Sender, input.TxHash)
	if err != nil {
		return nil, errors.Join(ErrGetBackupTransactionBySenderAndTxDoubleHashFail, err)
	}

	return &getBackupTransactionByHash.UCGetBackupTransactionByHash{
		ID:              transaction.ID,
		Sender:          transaction.Sender,
		Signature:       transaction.Signature,
		BlockNumber:     uint64(transaction.BlockNumber),
		EncryptedTx:     transaction.EncryptedTx,
		EncodingVersion: uint32(transaction.EncodingVersion),
		CreatedAt:       transaction.CreatedAt,
	}, nil
}
