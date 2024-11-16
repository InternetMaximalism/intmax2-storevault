package post_backup_user_state

import (
	"context"
	"errors"
	"intmax2-store-vault/configs"
	bpsTypes "intmax2-store-vault/internal/balance_prover_service/types"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	intMaxTypes "intmax2-store-vault/internal/types"
	postBackupUserState "intmax2-store-vault/internal/use_cases/post_backup_user_state"
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
) postBackupUserState.UseCasePostBackupUserState {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context, input *postBackupUserState.UCPostBackupUserStateInput,
) (*postBackupUserState.UCPostBackupUserState, error) {
	const (
		hName           = "UseCase PostBackupUserState"
		userAddressKey  = "user_address"
		balanceProofKey = "balance_proof"
		blockNumberKey  = "block_number"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCPostBackupUserStateInputEmpty)
		return nil, ErrUCPostBackupUserStateInputEmpty
	}

	span.SetAttributes(
		attribute.String(userAddressKey, input.UserAddress),
		attribute.String(balanceProofKey, input.BalanceProof),
		attribute.Int64(blockNumberKey, input.BlockNumber),
	)

	bp, err := intMaxTypes.NewCompressedPlonky2ProofFromBase64String(input.BalanceProof)
	if err != nil {
		return nil, errors.Join(ErrCompressedPlonky2ProofFromBase64StringFail, err)
	}

	var bpi *bpsTypes.BalancePublicInputs
	bpi, err = new(bpsTypes.BalancePublicInputs).FromPublicInputs(bp.PublicInputs)
	if err != nil {
		return nil, errors.Join(ErrBalancePublicInputsFromPublicInputsFail, err)
	}

	var bytesBP []byte
	bytesBP, err = bp.MarshalJSON()
	if err != nil {
		return nil, errors.Join(ErrMarshalPlonky2ProofWithBalanceProofFail, err)
	}

	var us *mDBApp.UserState
	us, err = u.db.CreateBackupUserState(
		input.UserAddress,
		input.EncryptedUserState,
		input.AuthSignature,
		input.BlockNumber,
	)
	if err != nil {
		return nil, errors.Join(ErrCreateBackupUserStateFail, err)
	}

	var bpDB *mDBApp.BalanceProof
	bpDB, err = u.db.CreateBalanceProof(
		us.ID, input.UserAddress, bpi.PrivateCommitment.String(), input.BlockNumber, bytesBP,
	)
	if err != nil {
		return nil, errors.Join(ErrCreateBalanceProofFail, err)
	}

	err = bp.UnmarshalJSON(bpDB.BalanceProof)
	if err != nil {
		return nil, errors.Join(ErrUnmarshalPlonky2ProofWithBalanceProofFail, err)
	}

	return &postBackupUserState.UCPostBackupUserState{
		ID:                 us.ID,
		UserAddress:        us.UserAddress,
		BalanceProof:       bp.ProofBase64String(),
		EncryptedUserState: us.EncryptedUserState,
		AuthSignature:      us.AuthSignature,
		BlockNumber:        us.BlockNumber,
		CreatedAt:          us.CreatedAt,
	}, nil
}
