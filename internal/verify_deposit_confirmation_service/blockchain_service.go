package verify_deposit_confirmation_service

import (
	"context"
)

//go:generate mockgen -destination=mock_blockchain_service_test.go -package=verify_deposit_confirmation_service_test -source=blockchain_service.go

type ServiceBlockchain interface {
	ChainSB
}

type ChainSB interface {
	EthereumNetworkChainLinkEvmJSONRPC(ctx context.Context) (string, error)
	ScrollNetworkChainLinkEvmJSONRPC(ctx context.Context) (string, error)
}
