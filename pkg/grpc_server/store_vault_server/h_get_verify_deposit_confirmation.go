package store_vault_server

import (
	"context"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	getVerifyDepositConfirmation "intmax2-store-vault/internal/use_cases/get_verify_deposit_confirmation"
	"intmax2-store-vault/pkg/grpc_server/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *StoreVaultServer) GetVerifyDepositConfirmation(
	ctx context.Context,
	req *node.GetVerifyDepositConfirmationRequest,
) (*node.GetVerifyDepositConfirmationResponse, error) {
	resp := node.GetVerifyDepositConfirmationResponse{}

	const (
		hName      = "Handler GetVerifyDepositConfirmation"
		requestKey = "request"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getVerifyDepositConfirmation.UCGetVerifyDepositConfirmationInput{
		DepositID: req.DepositId,
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var result *getVerifyDepositConfirmation.UCGetVerifyDepositConfirmation
	result, err = s.commands.GetVerifyDepositConfirmation(s.config, s.log, s.vdcs).Do(spanCtx, &input)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		const msg = "failed to get verify deposit confirmation request: %v"
		return &resp, utils.Internal(spanCtx, s.log, msg, err)
	}

	resp.Success = true
	resp.Data = &node.GetVerifyDepositConfirmationResponse_Data{Confirmed: result.IsVerifyDepositConfirmation}

	return &resp, utils.OK(spanCtx)
}
