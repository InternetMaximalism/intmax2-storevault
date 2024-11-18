package get_backup_deposits_list

import (
	mFL "intmax2-store-vault/internal/sql_filter/models"
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=get_backup_deposits_list_test -source=db_app.go

type SQLDriverApp interface {
	BackupDeposits
}

type BackupDeposits interface {
	GetBackupDepositsByRecipient(
		recipient string,
		pagination mDBApp.PaginationOfListOfBackupDepositsInput,
		sorting mFL.Sorting, orderBy mFL.OrderBy,
		filters mFL.FiltersList,
	) (
		paginator *mDBApp.PaginationOfListOfBackupDeposits,
		listDBApp mDBApp.ListOfBackupDeposit,
		err error,
	)
}
