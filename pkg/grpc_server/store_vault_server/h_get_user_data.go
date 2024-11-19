package store_vault_server

import (
	"context"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	"intmax2-store-vault/pkg/grpc_server/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *StoreVaultServer) GetUserData(
	ctx context.Context,
	req *node.GetUserDataRequest,
) (*node.GetUserDataResponse, error) {
	resp := node.GetUserDataResponse{}

	const (
		hName      = "Handler GetUserData"
		requestKey = "request"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	return &resp, utils.OK(spanCtx)
}
