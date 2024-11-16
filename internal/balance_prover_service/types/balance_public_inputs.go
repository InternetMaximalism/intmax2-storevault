package types

import (
	"errors"
	intMaxAcc "intmax2-store-vault/internal/accounts"
	bvpTypes "intmax2-store-vault/internal/block_validity_prover/types"
	intMaxGP "intmax2-store-vault/internal/hash/goldenposeidon"
	intMaxTypes "intmax2-store-vault/internal/types"

	"github.com/iden3/go-iden3-crypto/ffg"
)

type BalancePublicInputs struct {
	PubKey                  *intMaxAcc.PublicKey
	PrivateCommitment       *intMaxTypes.PoseidonHashOut
	LastTxHash              *intMaxTypes.PoseidonHashOut
	LastTxInsufficientFlags InsufficientFlags
	PublicState             *bvpTypes.PublicState
}

const (
	numHashOutElts                = intMaxGP.NUM_HASH_OUT_ELTS
	publicKeyOffset               = 0
	privateCommitmentOffset       = publicKeyOffset + int8Key
	lastTxHashOffset              = privateCommitmentOffset + numHashOutElts
	lastTxInsufficientFlagsOffset = lastTxHashOffset + numHashOutElts
	publicStateOffset             = lastTxInsufficientFlagsOffset + InsufficientFlagsLen
	sizeOfBalancePublicInputs     = publicStateOffset + bvpTypes.PublicStateLimbSize
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
	lastTxInsufficientFlags := new(InsufficientFlags).FromFieldElementSlice(
		publicInputs[lastTxInsufficientFlagsOffset:publicStateOffset],
	)
	publicState := new(bvpTypes.PublicState).FromFieldElementSlice(
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
