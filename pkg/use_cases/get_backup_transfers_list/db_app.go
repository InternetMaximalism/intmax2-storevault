package get_backup_transfers_list

import (
	mFL "intmax2-store-vault/internal/sql_filter/models"
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=get_backup_transfers_list_test -source=db_app.go

type SQLDriverApp interface {
	BackupTransfers
}

type BackupTransfers interface {
	GetBackupTransfersByRecipient(
		recipient string,
		pagination mDBApp.PaginationOfListOfBackupTransfersInput,
		sorting mFL.Sorting, orderBy mFL.OrderBy,
		filters mFL.FiltersList,
	) (
		paginator *mDBApp.PaginationOfListOfBackupTransfers,
		listDBApp mDBApp.ListOfBackupTransfer,
		err error,
	)
}
