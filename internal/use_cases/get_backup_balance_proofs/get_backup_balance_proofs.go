package get_backup_balance_proofs

import (
	"context"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_balance_proofs.go -package=mocks -source=get_backup_balance_proofs.go

type UCGetBackupBalanceProof struct {
	ID                         string
	EnoughBalanceProofBodyHash string
	LastBalanceProofBody       string
	BalanceTransitionProofBody string
}

type UCGetBackupBalanceProofs struct {
	Proofs []*UCGetBackupBalanceProof
}

type UCGetBackupBalanceProofsInput struct {
	Hashes []string `json:"hashes"`
}

// UseCaseGetBackupBalanceProofs describes GetBackupBalanceProofs contract.
type UseCaseGetBackupBalanceProofs interface {
	Do(ctx context.Context, input *UCGetBackupBalanceProofsInput) (*UCGetBackupBalanceProofs, error)
}
