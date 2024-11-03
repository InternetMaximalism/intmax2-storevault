package store_vault_service

import (
	intMaxAcc "intmax2-store-vault/internal/accounts"
	"intmax2-store-vault/internal/balance_prover_service"
	intMaxGP "intmax2-store-vault/internal/hash/goldenposeidon"
	"intmax2-store-vault/internal/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type DepositDetails struct {
	Recipient         *intMaxAcc.PublicKey
	TokenIndex        uint32
	Amount            *big.Int
	Salt              *intMaxGP.PoseidonHashOut
	RecipientSaltHash common.Hash
	DepositID         uint32
	DepositHash       common.Hash
}

type TransferDetails struct {
	TransferWitness              *types.TransferWitness
	TxTreeRoot                   intMaxGP.PoseidonHashOut
	TxIndex                      uint32
	TxMerkleProof                []intMaxGP.PoseidonHashOut
	SenderEnoughBalanceProofUUID string
}

type TransferDetailProof struct {
	SenderLastBalanceProof       *balance_prover_service.BalanceProofWithPublicInputs
	SenderBalanceTransitionProof *balance_prover_service.SpentProofWithPublicInputs
}

type TxDetails struct {
	Tx            *types.Tx
	TxTreeRoot    intMaxGP.PoseidonHashOut
	TxIndex       uint32
	TxMerkleProof []intMaxGP.PoseidonHashOut
}
