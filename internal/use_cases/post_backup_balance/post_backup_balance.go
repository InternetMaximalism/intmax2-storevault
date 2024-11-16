package post_backup_balance

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_post_backup_balance.go -package=mocks -source=post_backup_balance.go

type UCPostBackupBalance struct {
	ID                    string
	UserAddress           string
	EncryptedBalanceProof string
	EncryptedBalanceData  string
	EncryptedTxs          []string
	EncryptedTransfers    []string
	EncryptedDeposits     []string
	BlockNumber           uint64
	Signature             string
	CreatedAt             time.Time
}

type UCPostBackupBalanceInput struct {
	User                  string   `json:"user"`
	EncryptedBalanceProof string   `json:"encrypted_balance_proof"`
	EncryptedBalanceData  string   `json:"encrypted_balance_data"`
	EncryptedTxs          []string `json:"encrypted_txs"`
	EncryptedTransfers    []string `json:"encrypted_transfers"`
	EncryptedDeposits     []string `json:"encrypted_deposits"`
	Signature             string   `json:"signature"`
	BlockNumber           uint32   `json:"block_number"`
}

// UseCasePostBackupBalance describes PostBackupBalance contract.
type UseCasePostBackupBalance interface {
	Do(ctx context.Context, input *UCPostBackupBalanceInput) (*UCPostBackupBalance, error)
}
