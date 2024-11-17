package get_backup_transaction_by_uuid

import (
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=get_backup_transaction_by_uuid_test -source=db_app.go

type SQLDriverApp interface {
	BackupTransactions
}

type BackupTransactions interface {
	GetBackupTransactionByIDAndSender(id, sender string) (*mDBApp.BackupTransaction, error)
}
