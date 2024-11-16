package post_backup_user_state

import "errors"

// ErrUCPostBackupUserStateInputEmpty error: ucPostBackupUserStateInput must not be empty.
var ErrUCPostBackupUserStateInputEmpty = errors.New("ucPostBackupUserStateInput must not be empty")

// ErrCompressedPlonky2ProofFromBase64StringFail error: failed to get compressed Plonky2Proof from Base64String.
var ErrCompressedPlonky2ProofFromBase64StringFail = errors.New(
	"failed to get compressed Plonky2Proof from Base64String",
)

// ErrBalancePublicInputsFromPublicInputsFail error: failed to get BalancePublicInputs from PublicInputs.
var ErrBalancePublicInputsFromPublicInputsFail = errors.New("failed to get BalancePublicInputs from PublicInputs")

// ErrMarshalPlonky2ProofWithBalanceProofFail error: failed to marshal Plonky2Proof with BalanceProof.
var ErrMarshalPlonky2ProofWithBalanceProofFail = errors.New("failed to marshal Plonky2Proof with BalanceProof")

// ErrUnmarshalPlonky2ProofWithBalanceProofFail error: failed to unmarshal Plonky2Proof with BalanceProof.
var ErrUnmarshalPlonky2ProofWithBalanceProofFail = errors.New("failed to unmarshal Plonky2Proof with BalanceProof")

// ErrCreateBackupUserStateFail error: failed to create backup user state.
var ErrCreateBackupUserStateFail = errors.New("failed to create backup user state")

// ErrCreateBalanceProofFail error: failed to create balance proof.
var ErrCreateBalanceProofFail = errors.New("failed to create balance proof")
