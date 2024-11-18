package store_vault_server

import (
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	getBackupBalanceProofs "intmax2-store-vault/internal/use_cases/get_backup_balance_proofs"
	getBackupBalances "intmax2-store-vault/internal/use_cases/get_backup_balances"
	getBackupDepositByUuid "intmax2-store-vault/internal/use_cases/get_backup_deposit_by_uuid"
	getBackupDeposits "intmax2-store-vault/internal/use_cases/get_backup_deposits"
	getBackupDepositsList "intmax2-store-vault/internal/use_cases/get_backup_deposits_list"
	getBackupTransactionByUuid "intmax2-store-vault/internal/use_cases/get_backup_transaction_by_uuid"
	getBackupTransactions "intmax2-store-vault/internal/use_cases/get_backup_transactions"
	getBackupTransactionsList "intmax2-store-vault/internal/use_cases/get_backup_transactions_list"
	getBackupTransferByUuid "intmax2-store-vault/internal/use_cases/get_backup_transfer_by_uuid"
	getBackupTransfers "intmax2-store-vault/internal/use_cases/get_backup_transfers"
	getBackupTransfersList "intmax2-store-vault/internal/use_cases/get_backup_transfers_list"
	getBackupUserState "intmax2-store-vault/internal/use_cases/get_backup_user_state"
	getBalancesByAddress "intmax2-store-vault/internal/use_cases/get_balances_by_address"
	getVerifyDepositConfirmation "intmax2-store-vault/internal/use_cases/get_verify_deposit_confirmation"
	getVersion "intmax2-store-vault/internal/use_cases/get_version"
	postBackupBalance "intmax2-store-vault/internal/use_cases/post_backup_balance"
	postBackupDeposit "intmax2-store-vault/internal/use_cases/post_backup_deposit"
	postBackupTransaction "intmax2-store-vault/internal/use_cases/post_backup_transaction"
	postBackupTransfer "intmax2-store-vault/internal/use_cases/post_backup_transfer"
	postBackupUserState "intmax2-store-vault/internal/use_cases/post_backup_user_state"
	ucGetBackupBalanceProofs "intmax2-store-vault/pkg/use_cases/get_backup_balance_proofs"
	ucGetBackupBalances "intmax2-store-vault/pkg/use_cases/get_backup_balances"
	ucGetBackupDepositByUuid "intmax2-store-vault/pkg/use_cases/get_backup_deposit_by_uuid"
	ucGetBackupDeposits "intmax2-store-vault/pkg/use_cases/get_backup_deposits"
	ucGetBackupDepositsList "intmax2-store-vault/pkg/use_cases/get_backup_deposits_list"
	ucGetBackupTransactionByUuid "intmax2-store-vault/pkg/use_cases/get_backup_transaction_by_uuid"
	ucGetBackupTransactions "intmax2-store-vault/pkg/use_cases/get_backup_transactions"
	ucGetBackupTransactionsList "intmax2-store-vault/pkg/use_cases/get_backup_transactions_list"
	ucGetBackupTransferByUuid "intmax2-store-vault/pkg/use_cases/get_backup_transfer_by_uuid"
	ucGetBackupTransfers "intmax2-store-vault/pkg/use_cases/get_backup_transfers"
	ucGetBackupTransfersList "intmax2-store-vault/pkg/use_cases/get_backup_transfers_list"
	ucGetBackupUserState "intmax2-store-vault/pkg/use_cases/get_backup_user_state"
	ucGetBalancesByAddress "intmax2-store-vault/pkg/use_cases/get_balances_by_address"
	ucGetVerifyDepositConfirmation "intmax2-store-vault/pkg/use_cases/get_verify_deposit_confirmation"
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
	PostBackupBalance(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) postBackupBalance.UseCasePostBackupBalance
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
	GetBackupTransferByUuid(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupTransferByUuid.UseCaseGetBackupTransferByUuid
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
	GetBackupTransactionByUuid(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupTransactionByUuid.UseCaseGetBackupTransactionByUuid
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
	GetBackupDepositByUuid(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupDepositByUuid.UseCaseGetBackupDepositByUuid
	GetBackupBalances(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupBalances.UseCaseGetBackupBalances
	GetBalancesByAddress(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBalancesByAddress.UseCaseGetBalancesByAddress
	GetBackupSenderBalanceProofs(
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
	) getBackupBalanceProofs.UseCaseGetBackupBalanceProofs
	GetVerifyDepositConfirmation(
		cfg *configs.Config,
		log logger.Logger,
		vdcs VerifyDepositConfirmationService,
	) getVerifyDepositConfirmation.UseCaseGetVerifyDepositConfirmation
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

func (c *commands) PostBackupBalance(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) postBackupBalance.UseCasePostBackupBalance {
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

func (c *commands) GetBackupTransferByUuid(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupTransferByUuid.UseCaseGetBackupTransferByUuid {
	return ucGetBackupTransferByUuid.New(cfg, log, db)
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

func (c *commands) GetBackupTransactionByUuid(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupTransactionByUuid.UseCaseGetBackupTransactionByUuid {
	return ucGetBackupTransactionByUuid.New(cfg, log, db)
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

func (c *commands) GetBackupDepositByUuid(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupDepositByUuid.UseCaseGetBackupDepositByUuid {
	return ucGetBackupDepositByUuid.New(cfg, log, db)
}

func (c *commands) GetBackupBalances(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupBalances.UseCaseGetBackupBalances {
	return ucGetBackupBalances.New(cfg, log, db)
}

func (c *commands) GetBalancesByAddress(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBalancesByAddress.UseCaseGetBalancesByAddress {
	return ucGetBalancesByAddress.New(cfg, log, db)
}

func (c *commands) GetBackupSenderBalanceProofs(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupBalanceProofs.UseCaseGetBackupBalanceProofs {
	return ucGetBackupBalanceProofs.New(cfg, log, db)
}

func (c *commands) GetVerifyDepositConfirmation(
	cfg *configs.Config,
	log logger.Logger,
	vdcs VerifyDepositConfirmationService,
) getVerifyDepositConfirmation.UseCaseGetVerifyDepositConfirmation {
	return ucGetVerifyDepositConfirmation.New(cfg, log, vdcs)
}

func (c *commands) GetBackupUserState(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupUserState.UseCaseGetBackupUserState {
	return ucGetBackupUserState.New(cfg, log, db)
}
