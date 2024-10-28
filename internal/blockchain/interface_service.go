package blockchain

import (
	"context"
)

type ServiceBlockchain interface {
	ChainSB
}

type ChainSB interface {
	ScrollNetworkChainLinkEvmJSONRPC(ctx context.Context) (string, error)
	EthereumNetworkChainLinkEvmJSONRPC(ctx context.Context) (string, error)
}
