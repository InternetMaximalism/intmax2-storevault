package get_backup_user_state

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	intMaxTypes "intmax2-store-vault/internal/types"
	getBackupUserState "intmax2-store-vault/internal/use_cases/get_backup_user_state"
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"

	"go.opentelemetry.io/otel/attribute"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
}

func New(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupUserState.UseCaseGetBackupUserState {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context, input *getBackupUserState.UCGetBackupUserStateInput,
) (*getBackupUserState.UCGetBackupUserState, error) {
	const (
		hName          = "UseCase GetBackupUserState"
		userStateIDKey = "user_state_id"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupUserStateInputEmpty)
		return nil, ErrUCGetBackupUserStateInputEmpty
	}

	span.SetAttributes(
		attribute.String(userStateIDKey, input.UserStateID),
	)

	us, err := u.db.GetBackupUserState(input.UserStateID)
	if err != nil {
		return nil, errors.Join(ErrGetBackupUserStateFail, err)
	}

	var bpDB *mDBApp.BalanceProof
	bpDB, err = u.db.GetBalanceProofByUserStateID(us.ID)
	if err != nil {
		return nil, errors.Join(ErrGetBalanceProofByUserStateIDFail, err)
	}

	var bp intMaxTypes.Plonky2Proof
	err = bp.UnmarshalJSON(bpDB.BalanceProof)
	if err != nil {
		return nil, errors.Join(ErrUnmarshalPlonky2ProofWithBalanceProofFail, err)
	}

	return &getBackupUserState.UCGetBackupUserState{
		ID:                 us.ID,
		UserAddress:        us.UserAddress,
		BalanceProof:       bp.ProofBase64String(),
		EncryptedUserState: us.EncryptedUserState,
		AuthSignature:      us.AuthSignature,
		BlockNumber:        us.BlockNumber,
		CreatedAt:          us.CreatedAt,
	}, nil
}
