package get_backup_user_state

import "errors"

// ErrUCGetBackupUserStateInputEmpty error: ucGetBackupUserStateInput must not be empty.
var ErrUCGetBackupUserStateInputEmpty = errors.New("ucGetBackupUserStateInput must not be empty")

// ErrGetBackupUserStateFail error: failed to get backup user state.
var ErrGetBackupUserStateFail = errors.New("failed to get backup user state")

// ErrGetBalanceProofByUserStateIDFail error: failed to get balance proof by user state ID.
var ErrGetBalanceProofByUserStateIDFail = errors.New("failed to get balance proof by user state ID")

// ErrUnmarshalPlonky2ProofWithBalanceProofFail error: failed to unmarshal Plonky2Proof with BalanceProof.
var ErrUnmarshalPlonky2ProofWithBalanceProofFail = errors.New("failed to unmarshal Plonky2Proof with BalanceProof")
