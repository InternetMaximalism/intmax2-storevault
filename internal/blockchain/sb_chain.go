package blockchain

import (
	"context"
	"errors"
	"fmt"
	errorsB "intmax2-store-vault/internal/blockchain/errors"
	"intmax2-store-vault/internal/open_telemetry"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prodadidb/go-validation"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var ErrScrollChainIDInvalid = fmt.Errorf(errorsB.ErrScrollChainIDInvalidStr, ScrollMainNetChainID, ScrollSepoliaChainID)

var ErrEthereumChainIDInvalid = fmt.Errorf(
	errorsB.ErrEthereumChainIDInvalidStr, EthereumMainNetChainID, EthereumSepoliaChainID,
)

type ChainIDType string

const (
	EthereumMainNetChainID ChainIDType = "1"
	EthereumSepoliaChainID ChainIDType = "11155111"

	ScrollMainNetChainID ChainIDType = "534352"
	ScrollSepoliaChainID ChainIDType = "534351"
)

type ChainLinkEvmJSONRPC string

const (
	EthereumMainNetChainLinkEvmJSONRPC ChainLinkEvmJSONRPC = "https://mainnet.infura.io/v3"
	EthereumSepoliaChainLinkEvmJSONRPC ChainLinkEvmJSONRPC = "https://rpc.sepolia.org"

	ScrollMainNetChainLinkEvmJSONRPC ChainLinkEvmJSONRPC = "https://rpc.scroll.io"
	ScrollSepoliaChainLinkEvmJSONRPC ChainLinkEvmJSONRPC = "https://sepolia-rpc.scroll.io"
)

func (sb *serviceBlockchain) scrollNetworkChainIDValidator() error {
	return validation.Validate(sb.cfg.Blockchain.ScrollNetworkChainID,
		validation.Required,
		validation.In(
			string(ScrollMainNetChainID), string(ScrollSepoliaChainID),
		),
	)
}

func (sb *serviceBlockchain) ScrollNetworkChainLinkEvmJSONRPC(ctx context.Context) (string, error) {
	const (
		hName                   = "ServiceBlockchain func:ScrollNetworkChainLinkEvmJSONRPC"
		scrollNetworkChainIDKey = "scroll_network_chain_id"
		emptyKey                = ""
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(scrollNetworkChainIDKey, sb.cfg.Blockchain.ScrollNetworkChainID),
		))
	defer span.End()

	err := sb.scrollNetworkChainIDValidator()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return emptyKey, errors.Join(ErrScrollChainIDInvalid, err)
	}

	if strings.EqualFold(sb.cfg.Blockchain.ScrollNetworkChainID, string(ScrollMainNetChainID)) {
		return string(ScrollMainNetChainLinkEvmJSONRPC), nil
	}

	return string(ScrollSepoliaChainLinkEvmJSONRPC), nil
}

func (sb *serviceBlockchain) ethereumNetworkChainIDValidator() error {
	return validation.Validate(sb.cfg.Blockchain.EthereumNetworkChainID,
		validation.Required,
		validation.In(
			string(EthereumMainNetChainID), string(EthereumSepoliaChainID),
		),
	)
}

func (sb *serviceBlockchain) EthereumNetworkChainLinkEvmJSONRPC(ctx context.Context) (string, error) {
	const (
		hName                     = "ServiceBlockchain func:EthereumNetworkChainLinkEvmJSONRPC"
		ethereumNetworkChainIDKey = "ethereum_network_chain_id"
		emptyKey                  = ""
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(ethereumNetworkChainIDKey, sb.cfg.Blockchain.EthereumNetworkChainID),
		))
	defer span.End()

	sb.cfg.Blockchain.EthereumNetworkRpcUrl = strings.TrimSpace(sb.cfg.Blockchain.EthereumNetworkRpcUrl)
	if sb.cfg.Blockchain.EthereumNetworkRpcUrl != emptyKey {
		client, err := ethclient.DialContext(spanCtx, sb.cfg.Blockchain.EthereumNetworkRpcUrl)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			return emptyKey, errors.Join(errorsB.ErrEthClientDialFail)
		}

		var chainID *big.Int
		chainID, err = client.ChainID(spanCtx)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			return emptyKey, errors.Join(errorsB.ErrChainIDWithEthClientFail)
		}

		sb.cfg.Blockchain.EthereumNetworkChainID = chainID.String()

		return sb.cfg.Blockchain.EthereumNetworkRpcUrl, nil
	}

	err := sb.ethereumNetworkChainIDValidator()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return emptyKey, errors.Join(ErrEthereumChainIDInvalid, err)
	}

	if strings.EqualFold(sb.cfg.Blockchain.EthereumNetworkChainID, string(EthereumMainNetChainID)) {
		return string(EthereumMainNetChainLinkEvmJSONRPC), nil
	}

	return string(EthereumSepoliaChainLinkEvmJSONRPC), nil
}
