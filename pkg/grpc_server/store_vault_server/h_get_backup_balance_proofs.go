package store_vault_server

import (
	"context"
	"fmt"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	getBackupBalanceProofs "intmax2-store-vault/internal/use_cases/get_backup_balance_proofs"
	"intmax2-store-vault/pkg/grpc_server/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *StoreVaultServer) GetBackupBalanceProofs(ctx context.Context, req *node.GetBackupBalanceProofsRequest) (*node.GetBackupBalanceProofsResponse, error) {
	resp := node.GetBackupBalanceProofsResponse{}

	const (
		hName      = "Handler GetBackupBalanceProofs"
		requestKey = "request"
		actionKey  = "action"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := getBackupBalanceProofs.UCGetBackupBalanceProofsInput{
		Hashes: req.Hashes,
	}

	var info getBackupBalanceProofs.UCGetBackupBalanceProofs
	err := s.dbApp.Exec(spanCtx, &info, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var result *getBackupBalanceProofs.UCGetBackupBalanceProofs
		result, err = s.commands.GetBackupSenderBalanceProofs(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to get backup balances: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*getBackupBalanceProofs.UCGetBackupBalanceProofs); ok {
			v.Proofs = result.Proofs
		} else {
			const msg = "failed to convert of the backup balance proofs"
			err = fmt.Errorf(msg)
			open_telemetry.MarkSpanError(spanCtx, err)
			return err
		}

		return nil
	})
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		const msg = "failed to get the backup balance proofs by hashes"
		s.log.WithFields(logger.Fields{
			actionKey:  hName,
			requestKey: req.String(),
		}).WithError(err).Warnf(msg)

		return &resp, utils.OK(spanCtx)
	}

	resp.Success = true
	resp.Data = &node.GetBackupBalanceProofsResponse_Data{
		Proofs: make([]*node.GetBackupBalanceProofsResponse_Proof, len(info.Proofs)),
	}
	for key := range info.Proofs {
		resp.Data.Proofs[key] = &node.GetBackupBalanceProofsResponse_Proof{
			Id:                         info.Proofs[key].ID,
			EnoughBalanceProofBodyHash: info.Proofs[key].EnoughBalanceProofBodyHash,
			LastBalanceProofBody:       info.Proofs[key].LastBalanceProofBody,
			BalanceTransitionProofBody: info.Proofs[key].BalanceTransitionProofBody,
		}
	}

	return &resp, utils.OK(spanCtx)
}
