package get_backup_transfers

import (
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=get_backup_transfers_test -source=db_app.go

type SQLDriverApp interface {
	BackupTransfers
}

type BackupTransfers interface {
	GetBackupTransfers(condition string, value interface{}) ([]*mDBApp.BackupTransfer, error)
}
