package verify_deposit_confirmation_service

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/bindings"
	errorsB "intmax2-store-vault/internal/blockchain/errors"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	"intmax2-store-vault/pkg/utils"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.opentelemetry.io/otel/attribute"
)

type verifyDepositConfirmationService struct {
	ctx          context.Context
	cfg          *configs.Config
	log          logger.Logger
	client       *ethclient.Client
	scrollClient *ethclient.Client
	liquidity    *bindings.Liquidity
	rollup       *bindings.Rollup
}

func New(
	ctx context.Context,
	cfg *configs.Config,
	log logger.Logger,
	sb ServiceBlockchain,
) (VerifyDepositConfirmationService, error) {
	const (
		hName = "VerifyDepositConfirmationService func:New"

		moduleKey     = "module"
		moduleNameKey = "verify-deposit-confirmation-service"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	ethereumLink, err := sb.EthereumNetworkChainLinkEvmJSONRPC(ctx)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return nil, errors.Join(ErrScrollNetworkChainLinkEvmJSONRPCFail, err)
	}

	var ethereumClient *ethclient.Client
	ethereumClient, err = utils.NewClient(ethereumLink)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return nil, errors.Join(ErrCreateEthereumClientFail, err)
	}
	defer ethereumClient.Close()

	var liquidity *bindings.Liquidity
	liquidity, err = bindings.NewLiquidity(
		common.HexToAddress(cfg.Blockchain.LiquidityContractAddress),
		ethereumClient,
	)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return nil, errors.Join(ErrNewLiquidityContractFail, err)
	}

	var scrollLink string
	scrollLink, err = sb.ScrollNetworkChainLinkEvmJSONRPC(ctx)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return nil, errors.Join(ErrScrollNetworkChainLinkEvmJSONRPCFail, err)
	}

	var scrollClient *ethclient.Client
	scrollClient, err = utils.NewClient(scrollLink)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return nil, errors.Join(ErrCreateScrollClientFail, err)
	}
	defer scrollClient.Close()

	var rollup *bindings.Rollup
	rollup, err = bindings.NewRollup(
		common.HexToAddress(cfg.Blockchain.RollupContractAddress),
		scrollClient,
	)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return nil, errors.Join(ErrNewRollupContractFail, err)
	}

	return &verifyDepositConfirmationService{
		ctx:          ctx,
		cfg:          cfg,
		log:          log.WithFields(logger.Fields{moduleKey: moduleNameKey}),
		client:       ethereumClient,
		scrollClient: scrollClient,
		liquidity:    liquidity,
		rollup:       rollup,
	}, nil
}

func (v *verifyDepositConfirmationService) GetVerifyDepositConfirmation(
	ctx context.Context,
	depositID *big.Int,
) (bool, error) {
	const (
		hName        = "VerifyDepositConfirmationService func:GetVerifyDepositConfirmation"
		depositIDKey = "deposit_id"
		int1Key      = 1
		minus1Key    = -1
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if depositID == nil || depositID.Cmp(new(big.Int).SetUint64(int1Key)) == minus1Key {
		open_telemetry.MarkSpanError(spanCtx, ErrDepositIDInvalid)
		return false, ErrDepositIDInvalid
	}

	span.SetAttributes(attribute.String(depositIDKey, depositID.String()))

	lastProcessedDepositId, err := v.getLastProcessedDepositId(spanCtx)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return false, errors.Join(ErrGetLastProcessedDepositIdFail, err)
	}
	if lastProcessedDepositId.Cmp(depositID) == minus1Key {
		return false, nil
	}

	var depositExists bool
	depositExists, err = v.checkDepositDataExists(spanCtx, depositID)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return false, errors.Join(ErrCheckDepositDataExistsFail, err)
	}
	if !depositExists {
		return false, nil
	}

	var isDepositCanceled bool
	isDepositCanceled, err = v.checkIfDepositCanceled(spanCtx, depositID)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return false, errors.Join(ErrCheckIfDepositCanceledFail, err)
	}
	if isDepositCanceled {
		return false, nil
	}

	return true, nil
}

func (v *verifyDepositConfirmationService) getLastProcessedDepositId(
	ctx context.Context,
) (result *big.Int, err error) {
	const (
		hName                     = "VerifyDepositConfirmationService func:getLastProcessedDepositId"
		lastProcessedDepositIdKey = "last_processed_deposit_id"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	for {
		result, err = v.rollup.LastProcessedDepositId(&bind.CallOpts{
			Pending: false,
			Context: spanCtx,
		})
		if err != nil {
			if errorsB.ErrScrollProcessing(err, v.log, ErrLastProcessedDepositIdFromRollupContractFail.Error()) {
				<-time.After(time.Second)
				continue
			}

			open_telemetry.MarkSpanError(spanCtx, err)
			return nil, errors.Join(ErrLastProcessedDepositIdFromRollupContractFail, err)
		}

		break
	}

	span.SetAttributes(attribute.String(lastProcessedDepositIdKey, result.String()))

	return result, nil
}

func (v *verifyDepositConfirmationService) checkDepositDataExists(
	ctx context.Context,
	depositID *big.Int,
) (exists bool, err error) {
	const (
		hName                       = "VerifyDepositConfirmationService func:checkDepositDataExists"
		depositIDKey                = "deposit_id"
		isCheckDepositDataExistsKey = "is_check_deposit_data_exists"
		int1Key                     = 1
		minus1Key                   = -1
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if depositID == nil || depositID.Cmp(new(big.Int).SetUint64(int1Key)) == minus1Key {
		open_telemetry.MarkSpanError(spanCtx, ErrDepositIDInvalid)
		return false, ErrDepositIDInvalid
	}

	span.SetAttributes(attribute.String(depositIDKey, depositID.String()))

	var result bindings.DepositQueueLibDepositData
	for {
		result, err = v.liquidity.GetDepositData(&bind.CallOpts{
			Pending: false,
			Context: spanCtx,
		}, depositID)
		if err != nil {
			switch {
			case errorsB.ErrEthereumProcessing(err, v.log, ErrGetDepositDataFromLiquidityContractFail.Error()):
				<-time.After(time.Second)
				continue
			case strings.Contains(err.Error(), ErrOutOfBoundsAccessOfAnArrayOrBytesN):
				return false, nil
			default:
				open_telemetry.MarkSpanError(spanCtx, err)
				return false, errors.Join(ErrGetDepositDataFromLiquidityContractFail, err)
			}
		}

		break
	}

	exists = !result.IsRejected

	span.SetAttributes(attribute.Bool(isCheckDepositDataExistsKey, exists))

	return exists, nil
}

func (v *verifyDepositConfirmationService) checkIfDepositCanceled(
	ctx context.Context,
	depositID *big.Int,
) (bool, error) {
	const (
		hName                       = "VerifyDepositConfirmationService func:checkIfDepositCanceled"
		depositIDKey                = "deposit_id"
		isCheckIfDepositCanceledKey = "is_check_if_deposit_canceled"
		int1Key                     = 1
		minus1Key                   = -1
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if depositID == nil || depositID.Cmp(new(big.Int).SetUint64(int1Key)) == minus1Key {
		open_telemetry.MarkSpanError(spanCtx, ErrDepositIDInvalid)
		return false, ErrDepositIDInvalid
	}

	span.SetAttributes(attribute.String(depositIDKey, depositID.String()))

	depositIds := []*big.Int{depositID}
	iterator, err := v.liquidity.FilterDepositCanceled(&bind.FilterOpts{
		Start: v.cfg.Blockchain.LiquidityContractDeployedBlockNumber,
		End:   nil,
	}, depositIds)
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return false, errors.Join(ErrApplyFilterDepositCanceledWithLiquidityContractFail, err)
	}
	defer func() {
		_ = iterator.Close()
	}()

	isCanceled := false
	for iterator.Next() {
		if isCanceled {
			break
		}
		if iterator.Error() != nil {
			open_telemetry.MarkSpanError(spanCtx, iterator.Error())
			return false, errors.Join(ErrIteratingResultOfFilterDepositCanceledWithLiquidityContract, iterator.Error())
		}
		isCanceled = true
	}

	span.SetAttributes(attribute.Bool(isCheckIfDepositCanceledKey, isCanceled))

	return isCanceled, nil
}
