package backup_balance

import (
	"context"
	"encoding/binary"
	intMaxAcc "intmax2-store-vault/internal/accounts"
	"intmax2-store-vault/internal/block_validity_prover"
	"intmax2-store-vault/internal/finite_field"
	"intmax2-store-vault/internal/hash/goldenposeidon"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	intMaxTree "intmax2-store-vault/internal/tree"
	intMaxTypes "intmax2-store-vault/internal/types"
	"intmax2-store-vault/internal/use_cases/block_signature"
	"math/big"

	"github.com/iden3/go-iden3-crypto/ffg"
)

//go:generate mockgen -destination=../mocks/mock_post_backup_balance.go -package=mocks -source=post_backup_balance.go

const (
	numTransfersInTx     = 1 << intMaxTree.TRANSFER_TREE_HEIGHT
	InsufficientFlagsLen = numTransfersInTx / 32
	uint256LimbSize      = 8
	int4Key              = 4
	int32Key             = 32
)

type UCPostBackupBalance struct {
	Message string `json:"message"`
}

const (
	SuccessMsg = "Backup balance accepted."
)

type InsufficientFlags struct {
	Limbs [InsufficientFlagsLen]uint32
}

func (flags *InsufficientFlags) Equal(other *InsufficientFlags) bool {
	for i, limb := range flags.Limbs {
		if limb != other.Limbs[i] {
			return false
		}
	}

	return true
}

func (flags *InsufficientFlags) FromFieldElementSlice(value []ffg.Element) *InsufficientFlags {
	for i, x := range value {
		y := x.ToUint64Regular()
		if y >= uint64(1)<<int32Key {
			panic("overflow")
		}
		flags.Limbs[i] = uint32(y)
	}

	return flags
}

func (flags *InsufficientFlags) SetBit(index int, isValid bool) {
	limbIndex := index / int32Key
	bitIndex := index % int32Key

	if isValid {
		flags.Limbs[limbIndex] |= 1 << bitIndex
	} else {
		flags.Limbs[limbIndex] &^= 1 << bitIndex
	}
}

func (flags *InsufficientFlags) RandomAccess(index int) bool {
	limbIndex := index / int32Key
	bitIndex := index % int32Key

	return flags.Limbs[limbIndex]&(1<<bitIndex) != 0
}

func (flags *InsufficientFlags) Bytes() []byte {
	buf := make([]byte, InsufficientFlagsLen*int4Key)
	for i, limb := range flags.Limbs {
		binary.BigEndian.PutUint32(buf[i*4:], limb)
	}

	return buf
}

type BalancePublicInputs struct {
	PublicKey               *big.Int                           `json:"pubkey"`
	PrivateCommitment       goldenposeidon.PoseidonHashOut     `json:"privateCommitment"`
	LastTxHash              goldenposeidon.PoseidonHashOut     `json:"lastTxHash"`
	LastTxInsufficientFlags InsufficientFlags                  `json:"lastTxInsufficientFlags"`
	PublicState             *block_validity_prover.PublicState `json:"publicState"`
}

func (pis *BalancePublicInputs) Equal(other *BalancePublicInputs) bool {
	if pis.PublicKey.Cmp(other.PublicKey) != 0 {
		return false
	}
	if !pis.PrivateCommitment.Equal(&other.PrivateCommitment) {
		return false
	}
	if !pis.LastTxHash.Equal(&other.LastTxHash) {
		return false
	}
	if pis.LastTxInsufficientFlags != other.LastTxInsufficientFlags {
		return false
	}
	if !pis.PublicState.Equal(other.PublicState) {
		return false
	}
	return true
}

func VerifyEnoughBalanceProof(enoughBalanceProof *block_signature.Plonky2Proof) (*BalancePublicInputs, error) {
	publicInputs := make([]ffg.Element, len(enoughBalanceProof.PublicInputs))
	for i, publicInput := range enoughBalanceProof.PublicInputs {
		publicInputs[i].SetUint64(publicInput)
	}
	decodedPublicInputs := new(BalancePublicInputs).FromPublicInputs(publicInputs)
	err := decodedPublicInputs.Verify()
	if err != nil {
		return nil, err
	}

	// TODO: Verify verifier data in public inputs.

	// TODO: Verify enough balance proof by using Balance Validity Prover.
	return decodedPublicInputs, nil
}

func (pis *BalancePublicInputs) FromPublicInputs(publicInputs []ffg.Element) *BalancePublicInputs {
	const startPrivateCommitmentIndex = uint256LimbSize
	const lastTxHashIndex = startPrivateCommitmentIndex + goldenposeidon.NUM_HASH_OUT_ELTS
	const lastTxInsufficientFlagsIndex = lastTxHashIndex + goldenposeidon.NUM_HASH_OUT_ELTS
	const publicStateIndex = lastTxInsufficientFlagsIndex + InsufficientFlagsLen
	const endIndex = publicStateIndex + block_validity_prover.PublicStateLimbSize
	if len(publicInputs) != endIndex {
		panic("Invalid public inputs length")
	}

	publicKey := new(intMaxTypes.Uint256).FromFieldElementSlice(publicInputs[0:startPrivateCommitmentIndex])
	privateCommitment := new(goldenposeidon.PoseidonHashOut)
	copy(privateCommitment.Elements[:], publicInputs[startPrivateCommitmentIndex:lastTxHashIndex])
	lastTxHash := new(goldenposeidon.PoseidonHashOut)
	copy(lastTxHash.Elements[:], publicInputs[lastTxHashIndex:lastTxInsufficientFlagsIndex])
	lastTxInsufficientFlags := new(InsufficientFlags).FromFieldElementSlice(publicInputs[lastTxInsufficientFlagsIndex:publicStateIndex])
	publicState := new(block_validity_prover.PublicState).FromFieldElementSlice(publicInputs[publicStateIndex:endIndex])

	return &BalancePublicInputs{
		PublicKey:               publicKey.BigInt(),
		PrivateCommitment:       *privateCommitment,
		LastTxHash:              *lastTxHash,
		LastTxInsufficientFlags: *lastTxInsufficientFlags,
		PublicState:             publicState,
	}
}

func (pis *BalancePublicInputs) Verify() error {
	return nil
}

type EncryptedPlonky2Proof struct {
	Proof                 string `json:"proof"`
	EncryptedPublicInputs string `json:"publicInputs"`
}

type UCPostBackupBalanceInput struct {
	User                  string `json:"user"`
	EncryptedBalanceProof string `json:"encrypted_balance_proof"`
	EncryptedBalanceData  string `json:"encrypted_balance_data"`

	// List of transaction hashes already reflected
	EncryptedTxs []string `json:"encrypted_txs"`

	// List of transfer hashes already reflected
	EncryptedTransfers []string `json:"encrypted_transfers"`

	// List of deposit hashes already reflected
	EncryptedDeposits []string `json:"encrypted_deposits"`

	Signature string `json:"signature"`

	// DecodeUser            *intMaxAcc.PublicKey  `json:"-"`
	BlockNumber uint32 `json:"block_number"`
}

// UseCasePostBackupBalance describes PostBackupBalance contract.
type UseCasePostBackupBalance interface {
	Do(ctx context.Context, input *UCPostBackupBalanceInput) (*node.BackupBalanceResponse_Data_Balance, error)
}

func MakeMessage(
	user intMaxAcc.Address,
	blockNumber uint32,
	balanceProof []byte,
	encryptedBalancePublicInputs []byte,
	encryptedBalanceData []byte,
	encryptedTxs [][]byte,
	encryptedTransfers [][]byte,
	encryptedDeposits [][]byte,
) []ffg.Element {
	const numAddressBytes = 32
	buf := finite_field.NewBuffer(make([]ffg.Element, 0))
	finite_field.WriteFixedSizeBytes(buf, user.Bytes(), numAddressBytes)
	err := finite_field.WriteUint64(buf, uint64(blockNumber))
	// blockNumber is uint32, so it should be safe to cast to uint64
	if err != nil {
		panic(err)
	}
	finite_field.WriteBytes(buf, balanceProof)
	finite_field.WriteBytes(buf, encryptedBalancePublicInputs)
	finite_field.WriteBytes(buf, encryptedBalanceData)

	err = finite_field.WriteUint64(buf, uint64(len(encryptedTxs)))
	if err != nil {
		panic(err)
	}
	for _, tx := range encryptedTxs {
		finite_field.WriteBytes(buf, tx)
	}
	err = finite_field.WriteUint64(buf, uint64(len(encryptedTransfers)))
	if err != nil {
		panic(err)
	}
	for _, transfer := range encryptedTransfers {
		finite_field.WriteBytes(buf, transfer)
	}
	err = finite_field.WriteUint64(buf, uint64(len(encryptedDeposits)))
	if err != nil {
		panic(err)
	}
	for _, deposit := range encryptedDeposits {
		finite_field.WriteBytes(buf, deposit)
	}

	return buf.Inner()
}
