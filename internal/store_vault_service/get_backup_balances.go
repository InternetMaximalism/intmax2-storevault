package store_vault_service

import (
	"context"
	"fmt"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	backupBalance "intmax2-store-vault/internal/use_cases/backup_balance"
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

func GetBackupBalances(
	ctx context.Context,
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
	input *backupBalance.UCGetBackupBalancesInput,
) ([]*mDBApp.BackupBalance, error) {
	balances, err := db.GetLatestBackupBalanceByUserAddress(input.Sender, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup balances from db: %w", err)
	}
	return balances, nil
}
