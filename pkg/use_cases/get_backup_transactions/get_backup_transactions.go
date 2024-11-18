package get_backup_transactions

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	getBackupTransactions "intmax2-store-vault/internal/use_cases/get_backup_transactions"

	"go.opentelemetry.io/otel/attribute"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
}

func New(cfg *configs.Config, log logger.Logger, db SQLDriverApp) getBackupTransactions.UseCaseGetBackupTransactions {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context, input *getBackupTransactions.UCGetBackupTransactionsInput,
) (*getBackupTransactions.UCGetBackupTransactions, error) {
	const (
		hName               = "UseCase GetBackupTransactions"
		senderKey           = "sender"
		startBlockNumberKey = "start_block_number"
		limitKey            = "limit"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupTransactionsInputEmpty)
		return nil, ErrUCGetBackupTransactionsInputEmpty
	}

	span.SetAttributes(
		attribute.String(senderKey, input.Sender),
		attribute.Int64(startBlockNumberKey, int64(input.StartBlockNumber)),
		attribute.Int64(limitKey, int64(input.Limit)),
	)

	transactions, err := u.db.GetBackupTransactions("sender", input.Sender)
	if err != nil {
		return nil, errors.Join(ErrGetBackupTransactionsFail, err)
	}

	result := getBackupTransactions.UCGetBackupTransactions{
		Transactions: make([]*getBackupTransactions.UCGeyBackupTransactionsItem, len(transactions)),
		Meta: &getBackupTransactions.UCGetBackupTransactionsMeta{
			StartBlockNumber: input.StartBlockNumber,
			EndBlockNumber:   0,
		},
	}

	for key := range transactions {
		result.Transactions[key] = &getBackupTransactions.UCGeyBackupTransactionsItem{
			Uuid:            transactions[key].ID,
			Sender:          transactions[key].Sender,
			Signature:       transactions[key].Signature,
			BlockNumber:     uint64(transactions[key].BlockNumber),
			EncryptedTx:     transactions[key].EncryptedTx,
			EncodingVersion: uint32(transactions[key].EncodingVersion),
			CreatedAt:       transactions[key].CreatedAt,
		}
	}

	return &result, nil
}
