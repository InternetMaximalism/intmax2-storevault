package get_backup_deposit_by_uuid

import (
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=get_backup_deposit_by_uuid_test -source=db_app.go

type SQLDriverApp interface {
	BackupDeposits
}

type BackupDeposits interface {
	GetBackupDepositByIDAndRecipient(
		id, recipient string,
	) (*mDBApp.BackupDeposit, error)
}
