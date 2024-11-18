package get_backup_transaction_by_uuid

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	getBackupTransactionByUuid "intmax2-store-vault/internal/use_cases/get_backup_transaction_by_uuid"

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
) getBackupTransactionByUuid.UseCaseGetBackupTransactionByUuid {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context,
	input *getBackupTransactionByUuid.UCGetBackupTransactionByUuidInput,
) (*getBackupTransactionByUuid.UCGetBackupTransactionByUuid, error) {
	const (
		hName   = "UseCase GetBackupTransactionByUuid"
		uuidKey = "uuid"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupTransactionByUuidInputEmpty)
		return nil, ErrUCGetBackupTransactionByUuidInputEmpty
	}

	span.SetAttributes(
		attribute.String(uuidKey, input.Uuid),
	)

	transaction, err := u.db.GetBackupTransactionByID(input.Uuid)
	if err != nil {
		return nil, errors.Join(ErrGetBackupTransactionByIDFail, err)
	}

	return &getBackupTransactionByUuid.UCGetBackupTransactionByUuid{
		Uuid:            transaction.ID,
		Sender:          transaction.Sender,
		Signature:       transaction.Signature,
		BlockNumber:     uint64(transaction.BlockNumber),
		EncryptedTx:     transaction.EncryptedTx,
		EncodingVersion: uint32(transaction.EncodingVersion),
		CreatedAt:       transaction.CreatedAt,
	}, nil
}
