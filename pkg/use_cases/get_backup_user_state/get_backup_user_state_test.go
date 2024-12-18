package get_backup_user_state_test

import (
	"context"
	"errors"
	"fmt"
	"intmax2-store-vault/configs"
	bpsTypes "intmax2-store-vault/internal/balance_prover_service/types"
	intMaxTypes "intmax2-store-vault/internal/types"
	getBackupUserState "intmax2-store-vault/internal/use_cases/get_backup_user_state"
	"intmax2-store-vault/pkg/logger"
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
	errorsDB "intmax2-store-vault/pkg/sql_db/errors"
	ucGetBackupUserState "intmax2-store-vault/pkg/use_cases/get_backup_user_state"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUCGetBackupUserState(t *testing.T) {
	const int3Key = 3
	assert.NoError(t, configs.LoadDotEnv(int3Key))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := configs.New()
	log := logger.New(cfg.LOG.Level, cfg.LOG.TimeFormat, cfg.LOG.JSON, cfg.LOG.IsLogLine)
	db := NewMockSQLDriverApp(ctrl)

	uc := ucGetBackupUserState.New(cfg, log, db)

	const (
		path1 = "../../../"
		path2 = "./"
	)

	dir := path1
	if _, err := os.ReadFile(dir + cfg.APP.PEMPathCACert); err != nil {
		dir = path2
	}

	balanceProof, err := intMaxTypes.MakeSamplePlonky2Proof(dir)
	assert.NoError(t, err)
	b64 := balanceProof.ProofBase64String()
	bp, err := intMaxTypes.NewCompressedPlonky2ProofFromBase64String(b64)
	assert.NoError(t, err)
	byteBP, err := bp.MarshalJSON()
	assert.NoError(t, err)
	pi, err := new(bpsTypes.BalancePublicInputs).FromPublicInputs(bp.PublicInputs)
	assert.NoError(t, err)

	cases := []struct {
		desc    string
		input   *getBackupUserState.UCGetBackupUserState
		prepare func(*getBackupUserState.UCGetBackupUserState)
		err     error
	}{
		{
			desc: fmt.Sprintf("Error: %s", ucGetBackupUserState.ErrUCGetBackupUserStateInputEmpty),
			err:  ucGetBackupUserState.ErrUCGetBackupUserStateInputEmpty,
		},
		{
			desc: fmt.Sprintf("Error: %s", errorsDB.ErrNotFound),
			input: &getBackupUserState.UCGetBackupUserState{
				ID: uuid.New().String(),
			},
			prepare: func(state *getBackupUserState.UCGetBackupUserState) {
				db.EXPECT().GetBackupUserState(state.ID).Return(nil, errorsDB.ErrNotFound)
			},
			err: errorsDB.ErrNotFound,
		},
		{
			desc: fmt.Sprintf("Error: %s", errorsDB.ErrNotFound),
			input: &getBackupUserState.UCGetBackupUserState{
				ID: uuid.New().String(),
			},
			prepare: func(state *getBackupUserState.UCGetBackupUserState) {
				db.EXPECT().GetBackupUserState(state.ID).Return(&mDBApp.UserState{ID: state.ID}, nil)
				db.EXPECT().GetBalanceProofByUserStateID(state.ID).Return(nil, errorsDB.ErrNotFound)
			},
			err: errorsDB.ErrNotFound,
		},
		{
			desc: "Success",
			input: &getBackupUserState.UCGetBackupUserState{
				ID:                 uuid.New().String(),
				UserAddress:        uuid.New().String(),
				BalanceProof:       b64,
				EncryptedUserState: uuid.New().String(),
				AuthSignature:      uuid.New().String(),
				BlockNumber:        1,
				CreatedAt:          time.Now().UTC(),
			},
			prepare: func(state *getBackupUserState.UCGetBackupUserState) {
				db.EXPECT().GetBackupUserState(state.ID).Return(&mDBApp.UserState{
					ID:                 state.ID,
					UserAddress:        state.UserAddress,
					EncryptedUserState: state.EncryptedUserState,
					AuthSignature:      state.AuthSignature,
					BlockNumber:        state.BlockNumber,
					CreatedAt:          state.CreatedAt,
					UpdatedAt:          state.CreatedAt,
				}, nil)
				db.EXPECT().GetBalanceProofByUserStateID(state.ID).Return(&mDBApp.BalanceProof{
					ID:                     state.ID,
					UserStateID:            state.ID,
					UserAddress:            state.UserAddress,
					BlockNumber:            state.BlockNumber,
					PrivateStateCommitment: pi.PrivateCommitment.String(),
					BalanceProof:           byteBP,
					CreatedAt:              state.CreatedAt,
					UpdatedAt:              state.CreatedAt,
				}, nil)
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].desc, func(t *testing.T) {
			if cases[i].prepare != nil {
				cases[i].prepare(cases[i].input)
			}

			ctx := context.Background()
			var input *getBackupUserState.UCGetBackupUserStateInput
			if cases[i].input != nil {
				input = &getBackupUserState.UCGetBackupUserStateInput{
					UserStateID: cases[i].input.ID,
				}
			}

			var result *getBackupUserState.UCGetBackupUserState
			result, err = uc.Do(ctx, input)
			if cases[i].err != nil {
				assert.True(t, errors.Is(err, cases[i].err))
			} else {
				assert.NoError(t, err)
				if cases[i].input != nil {
					assert.Equal(t, result, cases[i].input)
				}
			}
		})
	}
}
