package get_backup_balances

import (
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=get_backup_balances_test -source=db_app.go

type SQLDriverApp interface {
	BackupBalances
}

type BackupBalances interface {
	GetLatestBackupBalanceByUserAddress(user string, limit int64) ([]*mDBApp.BackupBalance, error)
}
