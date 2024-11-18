package types

import (
	"encoding/binary"
	"errors"
	intMaxGP "intmax2-store-vault/internal/hash/goldenposeidon"
	intMaxTypes "intmax2-store-vault/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-crypto/ffg"
)

type PublicState struct {
	BlockTreeRoot       *intMaxGP.PoseidonHashOut `json:"blockTreeRoot"`
	PrevAccountTreeRoot *intMaxGP.PoseidonHashOut `json:"prevAccountTreeRoot"`
	AccountTreeRoot     *intMaxGP.PoseidonHashOut `json:"accountTreeRoot"`
	DepositTreeRoot     common.Hash               `json:"depositTreeRoot"`
	BlockHash           common.Hash               `json:"blockHash"`
	BlockNumber         uint32                    `json:"blockNumber"`
}

func (ps *PublicState) Equal(other *PublicState) bool {
	if !ps.BlockTreeRoot.Equal(other.BlockTreeRoot) {
		return false
	}
	if !ps.PrevAccountTreeRoot.Equal(other.PrevAccountTreeRoot) {
		return false
	}
	if !ps.AccountTreeRoot.Equal(other.AccountTreeRoot) {
		return false
	}
	if ps.DepositTreeRoot != other.DepositTreeRoot {
		return false
	}
	if ps.BlockHash != other.BlockHash {
		return false
	}
	if ps.BlockNumber != other.BlockNumber {
		return false
	}

	return true
}

const (
	prevAccountTreeRootOffset = intMaxGP.NUM_HASH_OUT_ELTS
	accountTreeRootOffset     = prevAccountTreeRootOffset + intMaxGP.NUM_HASH_OUT_ELTS
	depositTreeRootOffset     = accountTreeRootOffset + intMaxGP.NUM_HASH_OUT_ELTS
	blockHashOffset           = depositTreeRootOffset + int8Key
	blockNumberOffset         = blockHashOffset + int8Key
	PublicStateLimbSize       = blockNumberOffset + 1
)

func (ps *PublicState) FromFieldElementSlice(value []ffg.Element) *PublicState {
	ps.BlockTreeRoot = new(intMaxGP.PoseidonHashOut).FromPartial(value[:intMaxGP.NUM_HASH_OUT_ELTS])
	ps.PrevAccountTreeRoot = new(intMaxGP.PoseidonHashOut).FromPartial(value[prevAccountTreeRootOffset:accountTreeRootOffset])
	ps.AccountTreeRoot = new(intMaxGP.PoseidonHashOut).FromPartial(value[accountTreeRootOffset:depositTreeRootOffset])
	depositTreeRoot := intMaxTypes.Bytes32{}
	copy(depositTreeRoot[:], FieldElementSliceToUint32Slice(value[depositTreeRootOffset:blockHashOffset]))
	ps.DepositTreeRoot = common.Hash{}
	copy(ps.DepositTreeRoot[:], depositTreeRoot.Bytes())
	blockHash := intMaxTypes.Bytes32{}
	copy(blockHash[:], FieldElementSliceToUint32Slice(value[blockHashOffset:blockNumberOffset]))
	ps.BlockHash = common.Hash{}
	copy(ps.BlockHash[:], blockHash.Bytes())
	ps.BlockNumber = uint32(value[blockNumberOffset].ToUint64Regular())

	return ps
}

const NumPublicStateBytes = int32Key*int5Key + int4Key

func (ps *PublicState) Marshal() []byte {
	buf := make([]byte, NumPublicStateBytes)
	offset := 0

	copy(buf[offset:offset+int32Key], ps.BlockTreeRoot.Marshal())
	offset += int32Key

	copy(buf[offset:offset+int32Key], ps.PrevAccountTreeRoot.Marshal())
	offset += int32Key

	copy(buf[offset:offset+int32Key], ps.AccountTreeRoot.Marshal())
	offset += int32Key

	copy(buf[offset:offset+int32Key], ps.DepositTreeRoot.Bytes())
	offset += int32Key

	copy(buf[offset:offset+int32Key], ps.BlockHash.Bytes())

	binary.BigEndian.PutUint32(buf, ps.BlockNumber)

	return buf
}

func (ps *PublicState) Unmarshal(data []byte) error {
	if len(data) < NumPublicStateBytes {
		return errors.New("invalid data length")
	}

	offset := 0

	ps.BlockTreeRoot = new(intMaxGP.PoseidonHashOut)
	_ = ps.BlockTreeRoot.Unmarshal(data[offset : offset+int32Key])
	offset += int32Key

	ps.PrevAccountTreeRoot = new(intMaxGP.PoseidonHashOut)
	_ = ps.PrevAccountTreeRoot.Unmarshal(data[offset : offset+int32Key])
	offset += int32Key

	ps.AccountTreeRoot = new(intMaxGP.PoseidonHashOut)
	_ = ps.AccountTreeRoot.Unmarshal(data[offset : offset+int32Key])
	offset += int32Key

	ps.DepositTreeRoot = common.BytesToHash(data[offset : offset+int32Key])
	offset += int32Key

	ps.BlockHash = common.BytesToHash(data[offset : offset+int32Key])
	offset += int32Key

	ps.BlockNumber = binary.BigEndian.Uint32(data[offset : offset+int4Key])

	return nil
}
