package store_vault_server

import (
	"context"
	"math/big"
)

//go:generate mockgen -destination=mock_verify_deposit_confirmation_service_test.go -package=store_vault_server_test -source=verify_deposit_confirmation_service.go

type VerifyDepositConfirmationService interface {
	GetVerifyDepositConfirmation(
		ctx context.Context,
		depositID *big.Int,
	) (bool, error)
}
