package post_backup_user_state

import (
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=post_backup_user_state_test -source=db_app.go

type SQLDriverApp interface {
	BackupUserState
	BalanceProof
}

type BackupUserState interface {
	CreateBackupUserState(
		userAddress, encryptedUserState, authSignature string,
		blockNumber int64,
	) (*mDBApp.UserState, error)
}

type BalanceProof interface {
	CreateBalanceProof(
		userStateID, userAddress, privateStateCommitment string,
		blockNumber int64,
		balanceProof []byte,
	) (*mDBApp.BalanceProof, error)
}
