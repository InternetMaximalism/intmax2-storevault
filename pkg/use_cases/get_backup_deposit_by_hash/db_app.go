package get_backup_deposit_by_hash

import (
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=get_backup_deposit_by_hash_test -source=db_app.go

type SQLDriverApp interface {
	BackupDeposits
}

type BackupDeposits interface {
	GetBackupDepositByRecipientAndDepositDoubleHash(
		recipient, depositDoubleHash string,
	) (*mDBApp.BackupDeposit, error)
}
