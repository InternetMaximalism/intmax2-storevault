package get_verify_deposit_confirmation

import (
	"context"

	"github.com/holiman/uint256"
)

//go:generate mockgen -destination=../mocks/mock_get_verify_deposit_confirmation.go -package=mocks -source=get_verify_deposit_confirmation.go

type UCGetVerifyDepositConfirmation struct {
	IsVerifyDepositConfirmation bool
}

type UCGetVerifyDepositConfirmationInput struct {
	DepositID        string       `json:"deposit_id"`
	ConvertDepositID *uint256.Int `json:"-"`
}

// UseCaseGetVerifyDepositConfirmation describes GetVerifyDepositConfirmation contract.
type UseCaseGetVerifyDepositConfirmation interface {
	Do(ctx context.Context, input *UCGetVerifyDepositConfirmationInput) (*UCGetVerifyDepositConfirmation, error)
}
