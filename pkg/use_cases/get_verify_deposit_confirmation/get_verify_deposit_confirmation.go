package get_verify_deposit_confirmation

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	getVerifyDepositConfirmation "intmax2-store-vault/internal/use_cases/get_verify_deposit_confirmation"

	"go.opentelemetry.io/otel/attribute"
)

// uc describes use case
type uc struct {
	cfg  *configs.Config
	log  logger.Logger
	vdcs VerifyDepositConfirmationService
}

func New(
	cfg *configs.Config,
	log logger.Logger,
	vdcs VerifyDepositConfirmationService,
) getVerifyDepositConfirmation.UseCaseGetVerifyDepositConfirmation {
	return &uc{
		cfg:  cfg,
		log:  log,
		vdcs: vdcs,
	}
}

func (u *uc) Do(
	ctx context.Context,
	input *getVerifyDepositConfirmation.UCGetVerifyDepositConfirmationInput,
) (*getVerifyDepositConfirmation.UCGetVerifyDepositConfirmation, error) {
	const (
		hName     = "UseCase GetVerifyDepositConfirmation"
		depositId = "depositId"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetVerifyDepositConfirmationInputEmpty)
		return nil, ErrUCGetVerifyDepositConfirmationInputEmpty
	}

	span.SetAttributes(
		attribute.String(depositId, input.DepositID),
	)

	confirmed, err := u.vdcs.GetVerifyDepositConfirmation(spanCtx, input.ConvertDepositID.ToBig())
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return nil, errors.Join(ErrGetVerifyDepositConfirmationFail, err)
	}

	return &getVerifyDepositConfirmation.UCGetVerifyDepositConfirmation{
		IsVerifyDepositConfirmation: confirmed,
	}, nil
}
