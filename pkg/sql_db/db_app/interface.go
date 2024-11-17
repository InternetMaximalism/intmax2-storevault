package db_app

import (
	"context"
	mFL "intmax2-store-vault/internal/sql_filter/models"
	"intmax2-store-vault/pkg/sql_db/db_app/models"

	"github.com/dimiro1/health"
)

type SQLDb interface {
	GenericCommands
	ServiceCommands
	BackupTransfers
	BackupTransactions
	BackupDeposits
	BackupSenderProofs
	BackupBalances
	BackupUserState
	BalanceProof
}

type GenericCommands interface {
	Begin(ctx context.Context) (interface{}, error)
	Rollback()
	Commit() error
	Exec(ctx context.Context, input interface{}, executor func(d interface{}, input interface{}) error) (err error)
}

type ServiceCommands interface {
	Migrator(ctx context.Context, command string) (step int, err error)
	Check(ctx context.Context) health.Health
}

type BackupTransfers interface {
	CreateBackupTransfer(
		recipient, encryptedTransferHash, encryptedTransfer string,
		blockNumber int64,
	) (*models.BackupTransfer, error)
	GetBackupTransfer(condition string, value string) (*models.BackupTransfer, error)
	GetBackupTransferByIDAndRecipient(
		id, recipient string,
	) (*models.BackupTransfer, error)
	GetBackupTransfers(condition string, value interface{}) ([]*models.BackupTransfer, error)
	GetBackupTransfersByRecipient(
		recipient string,
		pagination models.PaginationOfListOfBackupTransfersInput,
		sorting mFL.Sorting, orderBy mFL.OrderBy,
		filters mFL.FiltersList,
	) (
		paginator *models.PaginationOfListOfBackupTransfers,
		listDBApp models.ListOfBackupTransfer,
		err error,
	)
}

type BackupTransactions interface {
	CreateBackupTransaction(
		sender, encryptedTxHash, encryptedTx, signature string,
		blockNumber int64,
	) (*models.BackupTransaction, error)
	GetBackupTransaction(condition string, value string) (*models.BackupTransaction, error)
	GetBackupTransactionByIDAndSender(id, sender string) (*models.BackupTransaction, error)
	GetBackupTransactions(condition string, value interface{}) ([]*models.BackupTransaction, error)
	GetBackupTransactionsBySender(
		sender string,
		pagination models.PaginationOfListOfBackupTransactionsInput,
		sorting mFL.Sorting, orderBy mFL.OrderBy,
		filters mFL.FiltersList,
	) (
		paginator *models.PaginationOfListOfBackupTransactions,
		listDBApp models.ListOfBackupTransaction,
		err error,
	)
}

type BackupDeposits interface {
	CreateBackupDeposit(
		recipient, depositHash, encryptedDeposit string,
		blockNumber int64,
	) (*models.BackupDeposit, error)
	GetBackupDepositByIDAndRecipient(
		id, recipient string,
	) (*models.BackupDeposit, error)
	GetBackupDeposit(conditions []string, values []interface{}) (*models.BackupDeposit, error)
	GetBackupDeposits(condition string, value interface{}) ([]*models.BackupDeposit, error)
	GetBackupDepositsByRecipient(
		recipient string,
		pagination models.PaginationOfListOfBackupDepositsInput,
		sorting mFL.Sorting, orderBy mFL.OrderBy,
		filters mFL.FiltersList,
	) (
		paginator *models.PaginationOfListOfBackupDeposits,
		listDBApp models.ListOfBackupDeposit,
		err error,
	)
}

type BackupSenderProofs interface {
	CreateBackupSenderProof(
		lastBalanceProofBody, balanceTransitionProofBody []byte,
		enoughBalanceProofBodyHash string,
	) (*models.BackupSenderProof, error)
	GetBackupSenderProofsByHashes(enoughBalanceProofBodyHashes []string) ([]*models.BackupSenderProof, error)
}

type BackupBalances interface {
	CreateBackupBalance(
		user, encryptedBalanceProof, encryptedBalanceData, signature string,
		encryptedTxs, encryptedTransfers, encryptedDeposits []string,
		blockNumber int64,
	) (*models.BackupBalance, error)
	GetBackupBalance(conditions []string, values []interface{}) (*models.BackupBalance, error)
	GetBackupBalances(condition string, value interface{}) ([]*models.BackupBalance, error)
	GetLatestBackupBalanceByUserAddress(user string, limit int64) ([]*models.BackupBalance, error)
}

type BackupUserState interface {
	CreateBackupUserState(
		userAddress, encryptedUserState, authSignature string,
		blockNumber int64,
	) (*models.UserState, error)
	UpdateBackupUserState(
		id, encryptedUserState, authSignature string,
		blockNumber int64,
	) (*models.UserState, error)
	GetBackupUserState(id string) (*models.UserState, error)
}

type BalanceProof interface {
	CreateBalanceProof(
		userStateID, userAddress, privateStateCommitment string,
		blockNumber int64,
		balanceProof []byte,
	) (*models.BalanceProof, error)
	GetBalanceProof(id string) (*models.BalanceProof, error)
	GetBalanceProofByUserStateID(userStateID string) (*models.BalanceProof, error)
}
