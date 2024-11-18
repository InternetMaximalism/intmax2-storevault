package blockchain

import (
	"context"
)

type ServiceBlockchain interface {
	ChainSB
}

type ChainSB interface {
	ScrollNetworkChainLinkEvmJSONRPC(ctx context.Context) (string, error)
	SetupEthereumNetworkChainID(ctx context.Context) error
	SetupScrollNetworkChainID(ctx context.Context) error
	EthereumNetworkChainLinkEvmJSONRPC(ctx context.Context) (string, error)
}
