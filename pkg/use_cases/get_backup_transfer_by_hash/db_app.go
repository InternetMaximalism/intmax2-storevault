package get_backup_transfer_by_hash

import (
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=get_backup_transfer_by_hash_test -source=db_app.go

type SQLDriverApp interface {
	BackupTransfers
}

type BackupTransfers interface {
	GetBackupTransferByRecipientAndTransferDoubleHash(
		recipient, transferDoubleHash string,
	) (*mDBApp.BackupTransfer, error)
}
