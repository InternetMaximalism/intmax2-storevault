package verify_deposit_confirmation_service

import (
	"errors"
)

const ErrOutOfBoundsAccessOfAnArrayOrBytesN = "execution reverted: out-of-bounds access of an array or bytesN"

// ErrCreateEthereumClientFail error: failed to create new ethereum client.
var ErrCreateEthereumClientFail = errors.New("failed to create new ethereum client")

// ErrScrollNetworkChainLinkEvmJSONRPCFail error: failed to get Scroll network chain link.
var ErrScrollNetworkChainLinkEvmJSONRPCFail = errors.New("failed to get Scroll network chain link")

// ErrCreateScrollClientFail error: failed to create new scroll client.
var ErrCreateScrollClientFail = errors.New("failed to create new scroll client")

// ErrNewLiquidityContractFail error: failed to instantiate a Liquidity contract.
var ErrNewLiquidityContractFail = errors.New("failed to instantiate a Liquidity contract")

// ErrNewRollupContractFail error: failed to instantiate a Rollup contract.
var ErrNewRollupContractFail = errors.New("failed to instantiate a Rollup contract")

// ErrDepositIDInvalid error: deposit ID must be valid.
var ErrDepositIDInvalid = errors.New("deposit ID must be valid")

// ErrLastProcessedDepositIdFromRollupContractFail error: failed to get last processed deposit ID from rollup contract.
var ErrLastProcessedDepositIdFromRollupContractFail = errors.New(
	"failed to get last processed deposit ID from rollup contract",
)

// ErrGetLastProcessedDepositIdFail error: failed to get last processed depositId.
var ErrGetLastProcessedDepositIdFail = errors.New("failed to get last processed depositId")

// ErrGetDepositDataFromLiquidityContractFail error: failed to get deposit data from liquidity contract.
var ErrGetDepositDataFromLiquidityContractFail = errors.New("failed to get deposit data from liquidity contract")

// ErrCheckDepositDataExistsFail error: failed to check of deposit data exists.
var ErrCheckDepositDataExistsFail = errors.New("failed to check of deposit data exists")

// ErrApplyFilterDepositCanceledWithLiquidityContractFail error: failed to apply filter deposit is canceled with liquidity contract.
var ErrApplyFilterDepositCanceledWithLiquidityContractFail = errors.New(
	"failed to apply filter deposit is canceled with liquidity contract",
)

// ErrIteratingResultOfFilterDepositCanceledWithLiquidityContract error: error encountered while iterating.
var ErrIteratingResultOfFilterDepositCanceledWithLiquidityContract = errors.New(
	"error encountered while iterating",
)

// ErrCheckIfDepositCanceledFail error: failed to check if deposit canceled.
var ErrCheckIfDepositCanceledFail = errors.New("failed to check if deposit canceled")
