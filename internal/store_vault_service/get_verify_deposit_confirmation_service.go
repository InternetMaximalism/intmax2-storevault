package store_vault_service

import (
	"context"
	"fmt"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/bindings"
	errorsB "intmax2-store-vault/internal/blockchain/errors"
	"intmax2-store-vault/internal/logger"
	verifyDepositConfirmation "intmax2-store-vault/internal/use_cases/verify_deposit_confirmation"
	"intmax2-store-vault/pkg/utils"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const Base10 = 10

type VerifyDepositConfirmationService struct {
	ctx          context.Context
	cfg          *configs.Config
	log          logger.Logger
	client       *ethclient.Client
	scrollClient *ethclient.Client
	liquidity    *bindings.Liquidity
	rollup       *bindings.Rollup
}

func newVerifyDepositConfirmationService(ctx context.Context, cfg *configs.Config, log logger.Logger, sb ServiceBlockchain) (*VerifyDepositConfirmationService, error) {
	client, err := utils.NewClient(cfg.Blockchain.EthereumNetworkRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}
	defer client.Close()

	scrollLink, err := sb.ScrollNetworkChainLinkEvmJSONRPC(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Scroll network chain link: %w", err)
	}

	scrollClient, err := utils.NewClient(scrollLink)
	if err != nil {
		return nil, fmt.Errorf("failed to create new scroll client: %w", err)
	}
	defer client.Close()

	liquidity, err := bindings.NewLiquidity(common.HexToAddress(cfg.Blockchain.LiquidityContractAddress), client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a Liquidity contract: %w", err)
	}

	rollup, err := bindings.NewRollup(common.HexToAddress(cfg.Blockchain.RollupContractAddress), scrollClient)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a Rollup contract: %w", err)
	}

	return &VerifyDepositConfirmationService{
		ctx:          ctx,
		cfg:          cfg,
		log:          log,
		client:       client,
		scrollClient: scrollClient,
		liquidity:    liquidity,
		rollup:       rollup,
	}, nil
}

func GetVerifyDepositConfirmation(
	ctx context.Context,
	cfg *configs.Config,
	log logger.Logger,
	sb ServiceBlockchain,
	input *verifyDepositConfirmation.UCGetVerifyDepositConfirmationInput,
) (bool, error) {
	service, err := newVerifyDepositConfirmationService(ctx, cfg, log, sb)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize VerifyDepositConfirmationService: %v", err.Error()))
	}

	depositId := new(big.Int)
	_, success := depositId.SetString(input.DepositId, Base10)
	if !success {
		panic(fmt.Sprintf("Failed to set depositId: %v", input.DepositId))
	}

	// TODO: Execute the following three tasks concurrently.
	lastProcessedDepositId, err := service.getLastProcessedDepositId()
	if err != nil {
		panic(fmt.Sprintf("Failed to get last processed depositId: %v", err.Error()))
	}
	if lastProcessedDepositId.Cmp(depositId) == -1 {
		return false, nil
	}

	depositExists, err := service.checkDepositDataExists(depositId)
	if err != nil {
		panic(fmt.Sprintf("Failed to check deposit data: %v", err.Error()))
	}
	if !depositExists {
		return false, nil
	}

	isDepositCanceled, err := service.checkIfDepositCanceled(depositId)
	if err != nil {
		panic(fmt.Sprintf("Failed to check deposit canceled: %v", err.Error()))
	}
	if isDepositCanceled {
		return false, nil
	}

	return true, nil
}

func (v *VerifyDepositConfirmationService) getLastProcessedDepositId() (result *big.Int, err error) {
	for {
		result, err = v.rollup.LastProcessedDepositId(&bind.CallOpts{
			Pending: false,
			Context: v.ctx,
		})
		if err != nil {
			if errorsB.ErrScrollProcessing(err, v.log, "failed to get last processed depositId from rollup") {
				<-time.After(time.Second)
				continue
			}

			return nil, fmt.Errorf("failed to get last processed depositId: %w", err)
		}

		break
	}

	return result, nil
}

func (v *VerifyDepositConfirmationService) checkDepositDataExists(depositId *big.Int) (exists bool, err error) {
	var result bindings.DepositQueueLibDepositData
	for {
		result, err = v.liquidity.GetDepositData(&bind.CallOpts{
			Pending: false,
			Context: v.ctx,
		}, depositId)
		if err != nil {
			switch {
			case errorsB.ErrEthereumProcessing(err, v.log, "failed to get deposit data from liquidity"):
				<-time.After(time.Second)
				continue
			case strings.Contains(err.Error(), "execution reverted: out-of-bounds access of an array or bytesN"):
				return false, nil
			default:
				return false, fmt.Errorf("failed to get deposit data: %w", err)
			}
		}

		break
	}

	exists = !result.IsRejected

	return exists, nil
}

func (v *VerifyDepositConfirmationService) checkIfDepositCanceled(depositId *big.Int) (bool, error) {
	depositIds := []*big.Int{depositId}
	iterator, err := v.liquidity.FilterDepositCanceled(&bind.FilterOpts{
		Start: 0,
		End:   nil,
	}, depositIds)
	if err != nil {
		return false, fmt.Errorf("failed to filter logs: %v", err)
	}

	defer iterator.Close()

	isCanceled := false
	for iterator.Next() {
		if iterator.Error() != nil {
			return false, fmt.Errorf("error encountered while iterating: %v", iterator.Error())
		}
		isCanceled = true
	}

	return isCanceled, nil
}
