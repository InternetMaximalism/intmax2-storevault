package configs

type Blockchain struct {
	ScrollNetworkChainID string `env:"BLOCKCHAIN_SCROLL_NETWORK_CHAIN_ID"`

	EthereumNetworkChainID string `env:"BLOCKCHAIN_ETHEREUM_NETWORK_CHAIN_ID"`
	EthereumNetworkRpcUrl  string `env:"BLOCKCHAIN_ETHEREUM_NETWORK_RPC_URL"`

	RollupContractAddress string `env:"BLOCKCHAIN_ROLLUP_CONTRACT_ADDRESS,required"`

	LiquidityContractAddress             string `env:"BLOCKCHAIN_LIQUIDITY_CONTRACT_ADDRESS,required"`
	LiquidityContractDeployedBlockNumber uint64 `env:"BLOCKCHAIN_LIQUIDITY_CONTRACT_DEPLOYED_BLOCK_NUMBER" envDefault:"0"`
}
