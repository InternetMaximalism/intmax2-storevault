package get_balances

import (
	"context"
	"fmt"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	service "intmax2-store-vault/internal/store_vault_service"
	backupBalance "intmax2-store-vault/internal/use_cases/backup_balance"

	"go.opentelemetry.io/otel/attribute"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
}

func New(cfg *configs.Config, log logger.Logger, db SQLDriverApp) backupBalance.UseCaseGetBalances {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context, input *backupBalance.UCGetBalancesInput,
) (*backupBalance.UCGetBalances, error) {
	const (
		hName   = "UseCase GetBalances"
		address = "address"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBalancesInputInputEmpty)
		return nil, ErrUCGetBalancesInputInputEmpty
	}

	span.SetAttributes(
		attribute.String(address, input.Address),
	)

	balances, err := service.GetBalances(ctx, u.cfg, u.log, u.db, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get balances: %w", err)
	}

	return balances, nil
}
