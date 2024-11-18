package get_backup_balance_proofs

import (
	"context"
	"encoding/base64"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	getBackupBalanceProofs "intmax2-store-vault/internal/use_cases/get_backup_balance_proofs"

	"go.opentelemetry.io/otel/attribute"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
}

func New(cfg *configs.Config, log logger.Logger, db SQLDriverApp) getBackupBalanceProofs.UseCaseGetBackupBalanceProofs {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context,
	input *getBackupBalanceProofs.UCGetBackupBalanceProofsInput,
) (*getBackupBalanceProofs.UCGetBackupBalanceProofs, error) {
	const (
		hName     = "UseCase GetBackupBalanceProofs"
		hashesKey = "hashes"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupBalanceProofsInputEmpty)
		return nil, ErrUCGetBackupBalanceProofsInputEmpty
	}

	span.SetAttributes(
		attribute.StringSlice(hashesKey, input.Hashes),
	)

	proofs, err := u.db.GetBackupSenderProofsByHashes(input.Hashes)
	if err != nil {
		var ErrGetBackupBalances = errors.New("failed to get backup balances")
		return nil, errors.Join(ErrGetBackupBalances, err)
	}

	results := getBackupBalanceProofs.UCGetBackupBalanceProofs{
		Proofs: make([]*getBackupBalanceProofs.UCGetBackupBalanceProof, len(proofs)),
	}
	for key := range proofs {
		results.Proofs[key] = &getBackupBalanceProofs.UCGetBackupBalanceProof{
			ID:                         proofs[key].ID,
			EnoughBalanceProofBodyHash: proofs[key].EnoughBalanceProofBodyHash,
			LastBalanceProofBody:       base64.StdEncoding.EncodeToString(proofs[key].LastBalanceProofBody),
			BalanceTransitionProofBody: base64.StdEncoding.EncodeToString(proofs[key].BalanceTransitionProofBody),
		}
	}

	return &results, nil
}
