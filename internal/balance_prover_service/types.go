package balance_prover_service

import (
	"errors"
	intMaxAcc "intmax2-store-vault/internal/accounts"
	"intmax2-store-vault/internal/block_validity_prover"
	intMaxGP "intmax2-store-vault/internal/hash/goldenposeidon"
	intMaxTypes "intmax2-store-vault/internal/types"
	"intmax2-store-vault/internal/use_cases/backup_balance"

	"github.com/iden3/go-iden3-crypto/ffg"
)

const (
	SENDER_TREE_HEIGHT     = 7
	balancePublicInputsLen = 47

	int2Key  = 2
	int3Key  = 3
	int4Key  = 4
	int8Key  = 8
	int32Key = 32
)

type poseidonHashOut = intMaxTypes.PoseidonHashOut

type BalancePublicInputs struct {
	PubKey                  *intMaxAcc.PublicKey
	PrivateCommitment       *intMaxTypes.PoseidonHashOut
	LastTxHash              *intMaxTypes.PoseidonHashOut
	LastTxInsufficientFlags backup_balance.InsufficientFlags
	PublicState             *block_validity_prover.PublicState
}

const (
	numHashOutElts                = intMaxGP.NUM_HASH_OUT_ELTS
	publicKeyOffset               = 0
	privateCommitmentOffset       = publicKeyOffset + int8Key
	lastTxHashOffset              = privateCommitmentOffset + numHashOutElts
	lastTxInsufficientFlagsOffset = lastTxHashOffset + numHashOutElts
	publicStateOffset             = lastTxInsufficientFlagsOffset + backup_balance.InsufficientFlagsLen
	sizeOfBalancePublicInputs     = publicStateOffset + block_validity_prover.PublicStateLimbSize
)

func (s *BalancePublicInputs) FromPublicInputs(publicInputs []ffg.Element) (*BalancePublicInputs, error) {
	if len(publicInputs) < balancePublicInputsLen {
		return nil, errors.New("invalid length")
	}

	address := new(intMaxTypes.Uint256).FromFieldElementSlice(publicInputs[0:int8Key])
	publicKey, err := new(intMaxAcc.PublicKey).SetBigInt(address.BigInt())
	if err != nil {
		return nil, err
	}
	privateCommitment := poseidonHashOut{
		Elements: [numHashOutElts]ffg.Element{
			publicInputs[privateCommitmentOffset],
			publicInputs[privateCommitmentOffset+1],
			publicInputs[privateCommitmentOffset+int2Key],
			publicInputs[privateCommitmentOffset+int3Key],
		},
	}
	lastTxHash := poseidonHashOut{
		Elements: [numHashOutElts]ffg.Element{
			publicInputs[lastTxHashOffset],
			publicInputs[lastTxHashOffset+1],
			publicInputs[lastTxHashOffset+int2Key],
			publicInputs[lastTxHashOffset+int3Key],
		},
	}
	lastTxInsufficientFlags := new(backup_balance.InsufficientFlags).FromFieldElementSlice(
		publicInputs[lastTxInsufficientFlagsOffset:publicStateOffset],
	)
	publicState := new(block_validity_prover.PublicState).FromFieldElementSlice(
		publicInputs[publicStateOffset:sizeOfBalancePublicInputs],
	)

	return &BalancePublicInputs{
		PubKey:                  publicKey,
		PrivateCommitment:       &privateCommitment,
		LastTxHash:              &lastTxHash,
		LastTxInsufficientFlags: *lastTxInsufficientFlags,
		PublicState:             publicState,
	}, nil
}
