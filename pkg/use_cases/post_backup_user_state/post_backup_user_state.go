package post_backup_user_state

import (
	"context"
	"fmt"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	service "intmax2-store-vault/internal/store_vault_service"
	postBackupUserState "intmax2-store-vault/internal/use_cases/post_backup_user_state"

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
) postBackupUserState.UseCasePostBackupUserState {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context, input *postBackupUserState.UCPostBackupUserStateInput,
) (*postBackupUserState.UCPostBackupUserState, error) {
	const (
		hName           = "UseCase PostBackupUserState"
		userAddressKey  = "user_address"
		balanceProofKey = "balance_proof"
		blockNumberKey  = "block_number"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCPostBackupUserStateInputEmpty)
		return nil, ErrUCPostBackupUserStateInputEmpty
	}

	span.SetAttributes(
		attribute.String(userAddressKey, input.UserAddress),
		attribute.String(balanceProofKey, input.BalanceProof),
		attribute.Int64(blockNumberKey, input.BlockNumber),
	)

	us, err := service.PostBackupUserState(ctx, u.cfg, u.log, u.db, input)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return nil, fmt.Errorf("failed to post backup user state: %w", err)
	}

	return us, nil
}
