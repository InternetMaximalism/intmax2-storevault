package get_backup_transfers

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	getBackupTransfers "intmax2-store-vault/internal/use_cases/get_backup_transfers"

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
) getBackupTransfers.UseCaseGetBackupTransfers {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context, input *getBackupTransfers.UCGetBackupTransfersInput,
) (*getBackupTransfers.UCGetBackupTransfers, error) {
	const (
		hName               = "UseCase GetBackupTransfers"
		senderKey           = "sender"
		startBlockNumberKey = "start_block_number"
		limitKey            = "limit"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupTransfersInputEmpty)
		return nil, ErrUCGetBackupTransfersInputEmpty
	}

	span.SetAttributes(
		attribute.String(senderKey, input.Sender),
		attribute.Int64(startBlockNumberKey, int64(input.StartBlockNumber)),
		attribute.Int64(limitKey, int64(input.Limit)),
	)

	transfers, err := u.db.GetBackupTransfers("recipient", input.Sender)
	if err != nil {
		return nil, errors.Join(ErrGetBackupTransfersFail, err)
	}

	result := getBackupTransfers.UCGetBackupTransfers{
		Transfers: make([]*getBackupTransfers.UCGetBackupTransfersTrItem, len(transfers)),
		Meta: &getBackupTransfers.UCGetBackupTransfersMeta{
			StartBlockNumber: input.StartBlockNumber,
			EndBlockNumber:   0,
		},
	}

	for key := range transfers {
		result.Transfers[key] = &getBackupTransfers.UCGetBackupTransfersTrItem{
			Uuid:              transfers[key].ID,
			BlockNumber:       transfers[key].BlockNumber,
			Recipient:         transfers[key].Recipient,
			EncryptedTransfer: transfers[key].EncryptedTransfer,
			CreatedAt:         transfers[key].CreatedAt,
		}
	}

	return &result, nil
}
