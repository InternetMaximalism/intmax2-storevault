package get_verify_deposit_confirmation

import (
	"context"
	"math/big"
)

//go:generate mockgen -destination=mock_verify_deposit_confirmation_service_test.go -package=get_verify_deposit_confirmation_test -source=verify_deposit_confirmation_service.go

type VerifyDepositConfirmationService interface {
	GetVerifyDepositConfirmation(
		ctx context.Context,
		depositID *big.Int,
	) (bool, error)
}
