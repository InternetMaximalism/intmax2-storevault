package get_backup_balances

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	getBackupBalances "intmax2-store-vault/internal/use_cases/get_backup_balances"

	"go.opentelemetry.io/otel/attribute"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
}

func New(cfg *configs.Config, log logger.Logger, db SQLDriverApp) getBackupBalances.UseCaseGetBackupBalances {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context, input *getBackupBalances.UCGetBackupBalancesInput,
) (*getBackupBalances.UCGetBackupBalances, error) {
	const (
		hName               = "UseCase GetBackupBalances"
		senderKey           = "sender"
		startBlockNumberKey = "start_block_number"
		limitKey            = "limit"
		int1Key             = 1
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupBalancesInputEmpty)
		return nil, ErrUCGetBackupBalancesInputEmpty
	}

	span.SetAttributes(
		attribute.String(senderKey, input.Sender),
		attribute.Int64(startBlockNumberKey, int64(input.StartBlockNumber)),
		attribute.Int64(limitKey, int64(input.Limit)),
	)

	balances, err := u.db.GetLatestBackupBalanceByUserAddress(input.Sender, int1Key)
	if err != nil {
		return nil, errors.Join(ErrGetLatestBackupBalanceByUserAddressFail, err)
	}

	result := getBackupBalances.UCGetBackupBalances{
		Balances: make([]*getBackupBalances.UCGetBackupBalancesItem, len(balances)),
		Meta: &getBackupBalances.UCGetBackupBalancesMeta{
			StartBlockNumber: input.StartBlockNumber,
			EndBlockNumber:   input.StartBlockNumber,
		},
	}

	for key := range balances {
		result.Balances[key] = &getBackupBalances.UCGetBackupBalancesItem{
			ID:                    balances[key].ID,
			UserAddress:           balances[key].UserAddress,
			EncryptedBalanceProof: balances[key].EncryptedBalanceProof,
			EncryptedBalanceData:  balances[key].EncryptedBalanceData,
			EncryptedTxs:          balances[key].EncryptedTxs,
			EncryptedTransfers:    balances[key].EncryptedTransfers,
			EncryptedDeposits:     balances[key].EncryptedDeposits,
			BlockNumber:           balances[key].BlockNumber,
			Signature:             balances[key].Signature,
			CreatedAt:             balances[key].CreatedAt,
		}
	}

	return &result, nil
}
