package get_backup_deposits

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	getBackupDeposits "intmax2-store-vault/internal/use_cases/get_backup_deposits"

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
) getBackupDeposits.UseCaseGetBackupDeposits {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context, input *getBackupDeposits.UCGetBackupDepositsInput,
) (*getBackupDeposits.UCGetBackupDeposits, error) {
	const (
		hName               = "UseCase GetBackupDeposits"
		senderKey           = "sender"
		startBlockNumberKey = "start_block_number"
		limitKey            = "limit"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupDepositsInputEmpty)
		return nil, ErrUCGetBackupDepositsInputEmpty
	}

	span.SetAttributes(
		attribute.String(senderKey, input.Sender),
		attribute.Int64(startBlockNumberKey, int64(input.StartBlockNumber)),
		attribute.Int64(limitKey, int64(input.Limit)),
	)

	deposits, err := u.db.GetBackupDeposits("recipient", input.Sender)
	if err != nil {
		return nil, errors.Join(ErrGetBackupDepositsFail, err)
	}

	result := getBackupDeposits.UCGetBackupDeposits{
		Deposits: make([]*getBackupDeposits.UCGetBackupDepositsItem, len(deposits)),
		Meta: &getBackupDeposits.UCGetBackupDepositsMeta{
			StartBlockNumber: input.StartBlockNumber,
			EndBlockNumber:   0,
		},
	}

	for key := range deposits {
		result.Deposits[key] = &getBackupDeposits.UCGetBackupDepositsItem{
			Uuid:             deposits[key].ID,
			Recipient:        deposits[key].Recipient,
			BlockNumber:      uint64(deposits[key].BlockNumber),
			EncryptedDeposit: deposits[key].EncryptedDeposit,
			CreatedAt:        deposits[key].CreatedAt,
		}
	}

	return &result, nil
}
