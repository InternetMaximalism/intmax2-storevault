package get_balances_by_address

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	getBackupBalancesByAddress "intmax2-store-vault/internal/use_cases/get_balances_by_address"
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"

	"go.opentelemetry.io/otel/attribute"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
}

func New(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupBalancesByAddress.UseCaseGetBalancesByAddress {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context,
	input *getBackupBalancesByAddress.UCGetBalancesByAddressInput,
) (*getBackupBalancesByAddress.UCGetBalancesByAddress, error) {
	const (
		hName   = "UseCase GetBalancesByAddress"
		address = "address"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBalancesByAddressInputEmpty)
		return nil, ErrUCGetBalancesByAddressInputEmpty
	}

	span.SetAttributes(
		attribute.String(address, input.Address),
	)

	// TODO: get these data concurrently

	deposits, err := u.db.GetBackupDeposits("recipient", input.Address)
	if err != nil {
		return nil, errors.Join(ErrGetBackupDepositsFail, err)
	}

	var transactions []*mDBApp.BackupTransaction
	transactions, err = u.db.GetBackupTransactions("sender", input.Address)
	if err != nil {
		return nil, errors.Join(ErrGetBackupTransactionsFail, err)
	}

	var transfers []*mDBApp.BackupTransfer
	transfers, err = u.db.GetBackupTransfers("recipient", input.Address)
	if err != nil {
		return nil, errors.Join(ErrGetBackupTransfersFail, err)
	}

	result := getBackupBalancesByAddress.UCGetBalancesByAddress{
		Deposits:     make([]*getBackupBalancesByAddress.GetBalancesByAddressBackupDeposit, len(deposits)),
		Transfers:    make([]*getBackupBalancesByAddress.GetBalancesByAddressBackupTransfer, len(transfers)),
		Transactions: make([]*getBackupBalancesByAddress.GetBalancesByAddressBackupTransaction, len(transactions)),
	}

	for key := range deposits {
		result.Deposits[key] = &getBackupBalancesByAddress.GetBalancesByAddressBackupDeposit{
			Recipient:        deposits[key].Recipient,
			EncryptedDeposit: deposits[key].EncryptedDeposit,
			BlockNumber:      uint64(deposits[key].BlockNumber),
			CreatedAt:        deposits[key].CreatedAt,
		}
	}

	for key := range transfers {
		result.Transfers[key] = &getBackupBalancesByAddress.GetBalancesByAddressBackupTransfer{
			EncryptedTransfer: transfers[key].EncryptedTransfer,
			Recipient:         transfers[key].Recipient,
			BlockNumber:       transfers[key].BlockNumber,
			CreatedAt:         transfers[key].CreatedAt,
		}
	}

	for key := range transactions {
		result.Transactions[key] = &getBackupBalancesByAddress.GetBalancesByAddressBackupTransaction{
			Sender:      transactions[key].Sender,
			EncryptedTx: transactions[key].EncryptedTx,
			BlockNumber: uint64(transactions[key].BlockNumber),
			CreatedAt:   transactions[key].CreatedAt,
		}
	}

	return &result, nil
}
