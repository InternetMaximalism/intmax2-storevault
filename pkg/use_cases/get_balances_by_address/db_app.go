package get_balances_by_address

import (
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=get_balances_by_address_test -source=db_app.go

type SQLDriverApp interface {
	BackupTransfers
	BackupTransactions
	BackupDeposits
}

type BackupTransfers interface {
	GetBackupTransfers(condition string, value interface{}) ([]*mDBApp.BackupTransfer, error)
}

type BackupTransactions interface {
	GetBackupTransactions(condition string, value interface{}) ([]*mDBApp.BackupTransaction, error)
}

type BackupDeposits interface {
	GetBackupDeposits(condition string, value interface{}) ([]*mDBApp.BackupDeposit, error)
}
