package store_vault_server

import (
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	backupBalance "intmax2-store-vault/internal/use_cases/backup_balance"
	backupBalanceProof "intmax2-store-vault/internal/use_cases/backup_balance_proof"
	getBackupDepositByHash "intmax2-store-vault/internal/use_cases/get_backup_deposit_by_hash"
	getBackupDeposits "intmax2-store-vault/internal/use_cases/get_backup_deposits"
	getBackupDepositsList "intmax2-store-vault/internal/use_cases/get_backup_deposits_list"
	getBackupTransactionByHash "intmax2-store-vault/internal/use_cases/get_backup_transaction_by_hash"
	getBackupTransactions "intmax2-store-vault/internal/use_cases/get_backup_transactions"
	getBackupTransactionsList "intmax2-store-vault/internal/use_cases/get_backup_transactions_list"
	getBackupTransferByHash "intmax2-store-vault/internal/use_cases/get_backup_transfer_by_hash"
	getBackupTransfers "intmax2-store-vault/internal/use_cases/get_backup_transfers"
	getBackupTransfersList "intmax2-store-vault/internal/use_cases/get_backup_transfers_list"
	getBackupUserState "intmax2-store-vault/internal/use_cases/get_backup_user_state"
	getVersion "intmax2-store-vault/internal/use_cases/get_version"
	postBackupDeposit "intmax2-store-vault/internal/use_cases/post_backup_deposit"
	postBackupTransaction "intmax2-store-vault/internal/use_cases/post_backup_transaction"
	postBackupTransfer "intmax2-store-vault/internal/use_cases/post_backup_transfer"
	postBackupUserState "intmax2-store-vault/internal/use_cases/post_backup_user_state"
	verifyDepositConfirmation "intmax2-store-vault/internal/use_cases/verify_deposit_confirmation"
	ucGetBackupBalanceProofs "intmax2-store-vault/pkg/use_cases/get_backup_balance_proofs"
	ucGetBackupBalances "intmax2-store-vault/pkg/use_cases/get_backup_balances"
	ucGetBackupDepositByHash "intmax2-store-vault/pkg/use_cases/get_backup_deposit_by_hash"
	ucGetBackupDeposits "intmax2-store-vault/pkg/use_cases/get_backup_deposits"
	ucGetBackupDepositsList "intmax2-store-vault/pkg/use_cases/get_backup_deposits_list"
	ucGetBackupTransactionByHash "intmax2-store-vault/pkg/use_cases/get_backup_transaction_by_hash"
	ucGetBackupTransactions "intmax2-store-vault/pkg/use_cases/get_backup_transactions"
	ucGetBackupTransactionsList "intmax2-store-vault/pkg/use_cases/get_backup_transactions_list"
	ucGetBackupTransferByHash "intmax2-store-vault/pkg/use_cases/get_backup_transfer_by_hash"
	ucGetBackupTransfers "intmax2-store-vault/pkg/use_cases/get_backup_transfers"
	ucGetBackupTransfersList "intmax2-store-vault/pkg/use_cases/get_backup_transfers_list"
	ucGetBackupUserState "intmax2-store-vault/pkg/use_cases/get_backup_user_state"
	ucGetBalances "intmax2-store-vault/pkg/use_cases/get_balances"
	ucVerifyDepositConfirmation "intmax2-store-vault/pkg/use_cases/get_verify_deposit_confirmation"
	ucGetVersion "intmax2-store-vault/pkg/use_cases/get_version"
	ucPostBackupBalance "intmax2-store-vault/pkg/use_cases/post_backup_balance"
	ucPostBackupDeposit "intmax2-store-vault/pkg/use_cases/post_backup_deposit"
	ucPostBackupTransaction "intmax2-store-vault/pkg/use_cases/post_backup_transaction"
	ucPostBackupTransfer "intmax2-store-vault/pkg/use_cases/post_backup_transfer"
	ucPostBackupUserState "intmax2-store-vault/pkg/use_cases/post_backup_user_state"
)

//go:generate mockgen -destination=mock_commands_test.go -package=store_vault_server_test -source=commands.go

type Commands interface {
	GetVersion(version, buildTime string) getVersion.UseCaseGetVersion
	PostBackupTransfer(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) postBackupTransfer.UseCasePostBackupTransfer
	PostBackupTransaction(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) postBackupTransaction.UseCasePostBackupTransaction
	PostBackupDeposit(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) postBackupDeposit.UseCasePostBackupDeposit
	PostBackupBalance(cfg *configs.Config, log logger.Logger, db SQLDriverApp) backupBalance.UseCasePostBackupBalance
	PostBackupUserState(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) postBackupUserState.UseCasePostBackupUserState
	GetBackupTransfers(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupTransfers.UseCaseGetBackupTransfers
	GetBackupTransfersList(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupTransfersList.UseCaseGetBackupTransfersList
	GetBackupTransferByHash(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupTransferByHash.UseCaseGetBackupTransferByHash
	GetBackupTransactions(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupTransactions.UseCaseGetBackupTransactions
	GetBackupTransactionsList(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupTransactionsList.UseCaseGetBackupTransactionsList
	GetBackupTransactionByHash(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupTransactionByHash.UseCaseGetBackupTransactionByHash
	GetBackupDeposits(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupDeposits.UseCaseGetBackupDeposits
	GetBackupDepositsList(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupDepositsList.UseCaseGetBackupDepositsList
	GetBackupDepositByHash(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupDepositByHash.UseCaseGetBackupDepositByHash
	GetBackupBalances(cfg *configs.Config, log logger.Logger, db SQLDriverApp) backupBalance.UseCaseGetBackupBalances
	GetBalances(cfg *configs.Config, log logger.Logger, db SQLDriverApp) backupBalance.UseCaseGetBalances
	GetBackupSenderBalanceProofs(cfg *configs.Config, log logger.Logger, db SQLDriverApp) backupBalanceProof.UseCaseGetBackupBalanceProofs
	GetVerifyDepositConfirmation(cfg *configs.Config, log logger.Logger, sb ServiceBlockchain) verifyDepositConfirmation.UseCaseGetVerifyDepositConfirmation
	GetBackupUserState(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupUserState.UseCaseGetBackupUserState
}

type commands struct{}

func NewCommands() Commands {
	return &commands{}
}

func (c *commands) GetVersion(version, buildTime string) getVersion.UseCaseGetVersion {
	return ucGetVersion.New(version, buildTime)
}

func (c *commands) PostBackupTransfer(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) postBackupTransfer.UseCasePostBackupTransfer {
	return ucPostBackupTransfer.New(cfg, log, db)
}

func (c *commands) PostBackupTransaction(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) postBackupTransaction.UseCasePostBackupTransaction {
	return ucPostBackupTransaction.New(cfg, log, db)
}

func (c *commands) PostBackupDeposit(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) postBackupDeposit.UseCasePostBackupDeposit {
	return ucPostBackupDeposit.New(cfg, log, db)
}

func (c *commands) PostBackupBalance(cfg *configs.Config, log logger.Logger, db SQLDriverApp) backupBalance.UseCasePostBackupBalance {
	return ucPostBackupBalance.New(cfg, log, db)
}

func (c *commands) PostBackupUserState(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) postBackupUserState.UseCasePostBackupUserState {
	return ucPostBackupUserState.New(cfg, log, db)
}

func (c *commands) GetBackupTransfers(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupTransfers.UseCaseGetBackupTransfers {
	return ucGetBackupTransfers.New(cfg, log, db)
}

func (c *commands) GetBackupTransfersList(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupTransfersList.UseCaseGetBackupTransfersList {
	return ucGetBackupTransfersList.New(cfg, log, db)
}

func (c *commands) GetBackupTransferByHash(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupTransferByHash.UseCaseGetBackupTransferByHash {
	return ucGetBackupTransferByHash.New(cfg, log, db)
}

func (c *commands) GetBackupTransactions(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupTransactions.UseCaseGetBackupTransactions {
	return ucGetBackupTransactions.New(cfg, log, db)
}

func (c *commands) GetBackupTransactionsList(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupTransactionsList.UseCaseGetBackupTransactionsList {
	return ucGetBackupTransactionsList.New(cfg, log, db)
}

func (c *commands) GetBackupTransactionByHash(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupTransactionByHash.UseCaseGetBackupTransactionByHash {
	return ucGetBackupTransactionByHash.New(cfg, log, db)
}

func (c *commands) GetBackupDeposits(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupDeposits.UseCaseGetBackupDeposits {
	return ucGetBackupDeposits.New(cfg, log, db)
}

func (c *commands) GetBackupDepositsList(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupDepositsList.UseCaseGetBackupDepositsList {
	return ucGetBackupDepositsList.New(cfg, log, db)
}

func (c *commands) GetBackupDepositByHash(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupDepositByHash.UseCaseGetBackupDepositByHash {
	return ucGetBackupDepositByHash.New(cfg, log, db)
}

func (c *commands) GetBackupBalances(cfg *configs.Config, log logger.Logger, db SQLDriverApp) backupBalance.UseCaseGetBackupBalances {
	return ucGetBackupBalances.New(cfg, log, db)
}

func (c *commands) GetBalances(cfg *configs.Config, log logger.Logger, db SQLDriverApp) backupBalance.UseCaseGetBalances {
	return ucGetBalances.New(cfg, log, db)
}

func (c *commands) GetBackupSenderBalanceProofs(cfg *configs.Config, log logger.Logger, db SQLDriverApp) backupBalanceProof.UseCaseGetBackupBalanceProofs {
	return ucGetBackupBalanceProofs.New(cfg, log, db)
}

func (c *commands) GetVerifyDepositConfirmation(cfg *configs.Config, log logger.Logger, sb ServiceBlockchain) verifyDepositConfirmation.UseCaseGetVerifyDepositConfirmation {
	return ucVerifyDepositConfirmation.New(cfg, log, sb)
}

func (c *commands) GetBackupUserState(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupUserState.UseCaseGetBackupUserState {
	return ucGetBackupUserState.New(cfg, log, db)
}
