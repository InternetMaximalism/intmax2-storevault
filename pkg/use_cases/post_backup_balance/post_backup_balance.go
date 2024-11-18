package post_backup_balance

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	postBackupBalance "intmax2-store-vault/internal/use_cases/post_backup_balance"

	"go.opentelemetry.io/otel/attribute"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
}

func New(cfg *configs.Config, log logger.Logger, db SQLDriverApp) postBackupBalance.UseCasePostBackupBalance {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context, input *postBackupBalance.UCPostBackupBalanceInput,
) (*postBackupBalance.UCPostBackupBalance, error) {
	const (
		hName                    = "UseCase PostBackupBalance"
		userKey                  = "user"
		blockNumberKey           = "block_number"
		encryptedBalanceProofKey = "encrypted_balance_proof"
		encryptedBalanceDataKey  = "encrypted_balance_data"
		encryptedTxsKey          = "encrypted_txs"
		encryptedTransfersKey    = "encrypted_transfers"
		encryptedDepositsKey     = "encrypted_deposits"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCPostBackupBalanceInputEmpty)
		return nil, ErrUCPostBackupBalanceInputEmpty
	}

	span.SetAttributes(
		attribute.String(userKey, input.User),
		attribute.String(encryptedBalanceProofKey, input.EncryptedBalanceProof),
		attribute.String(encryptedBalanceDataKey, input.EncryptedBalanceData),
		attribute.StringSlice(encryptedTxsKey, input.EncryptedTxs),
		attribute.StringSlice(encryptedTransfersKey, input.EncryptedTransfers),
		attribute.StringSlice(encryptedDepositsKey, input.EncryptedDeposits),
		attribute.Int64(blockNumberKey, int64(input.BlockNumber)),
	)

	newBackupBalance, err := u.db.CreateBackupBalance(
		input.User,
		input.EncryptedBalanceProof,
		input.EncryptedBalanceData,
		input.Signature,
		input.EncryptedTxs,
		input.EncryptedTransfers,
		input.EncryptedDeposits,
		int64(input.BlockNumber),
	)
	if err != nil {
		return nil, errors.Join(ErrCreateBackupBalanceFail, err)
	}

	return &postBackupBalance.UCPostBackupBalance{
		ID:                    newBackupBalance.ID,
		UserAddress:           newBackupBalance.UserAddress,
		EncryptedBalanceProof: newBackupBalance.EncryptedBalanceProof,
		EncryptedBalanceData:  newBackupBalance.EncryptedBalanceData,
		EncryptedTxs:          newBackupBalance.EncryptedTxs,
		EncryptedTransfers:    newBackupBalance.EncryptedTransfers,
		EncryptedDeposits:     newBackupBalance.EncryptedDeposits,
		BlockNumber:           newBackupBalance.BlockNumber,
		Signature:             newBackupBalance.Signature,
		CreatedAt:             newBackupBalance.CreatedAt,
	}, nil
}
