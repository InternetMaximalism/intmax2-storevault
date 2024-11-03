package store_vault_service

import (
	"context"
	"fmt"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	getBackupTransferByHash "intmax2-store-vault/internal/use_cases/get_backup_transfer_by_hash"
	getBackupTransfers "intmax2-store-vault/internal/use_cases/get_backup_transfers"
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
)

func GetBackupTransfers(
	ctx context.Context,
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
	input *getBackupTransfers.UCGetBackupTransfersInput,
) ([]*mDBApp.BackupTransfer, error) {
	transfers, err := db.GetBackupTransfers("recipient", input.Sender)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup transfers from db: %w", err)
	}
	fmt.Println("transfers", transfers)
	return transfers, nil
}

func GetBackupTransferByHash(
	ctx context.Context,
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
	input *getBackupTransferByHash.UCGetBackupTransferByHashInput,
) (*mDBApp.BackupTransfer, error) {
	transfer, err := db.GetBackupTransferByRecipientAndTransferDoubleHash(input.Recipient, input.TransferHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup transfer by hash from db: %w", err)
	}
	return transfer, nil
}
