package get_backup_balance_proofs

import (
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

//go:generate mockgen -destination=mock_db_app_test.go -package=get_backup_balance_proofs_test -source=db_app.go

type SQLDriverApp interface {
	BackupSenderProofs
}

type BackupSenderProofs interface {
	GetBackupSenderProofsByHashes(enoughBalanceProofBodyHashes []string) ([]*mDBApp.BackupSenderProof, error)
}
