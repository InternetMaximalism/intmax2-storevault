package get_backup_user_state

import (
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=get_backup_user_state_test -source=db_app.go

type SQLDriverApp interface {
	BackupUserState
	BalanceProof
}

type BackupUserState interface {
	GetBackupUserState(id string) (*mDBApp.UserState, error)
}

type BalanceProof interface {
	GetBalanceProofByUserStateID(userStateID string) (*mDBApp.BalanceProof, error)
}
