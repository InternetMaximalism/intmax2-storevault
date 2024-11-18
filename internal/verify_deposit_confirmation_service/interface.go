package verify_deposit_confirmation_service

import (
	"context"
	"math/big"
)

type VerifyDepositConfirmationService interface {
	GetVerifyDepositConfirmation(
		ctx context.Context,
		depositID *big.Int,
	) (bool, error)
}
