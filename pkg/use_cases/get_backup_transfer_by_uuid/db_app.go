package get_backup_transfer_by_uuid

import (
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=get_backup_transfer_by_uuid_test -source=db_app.go

type SQLDriverApp interface {
	BackupTransfers
}

type BackupTransfers interface {
	GetBackupTransferByIDAndRecipient(
		id, recipient string,
	) (*mDBApp.BackupTransfer, error)
}
