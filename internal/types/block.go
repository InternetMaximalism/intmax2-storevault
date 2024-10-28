package types

import "github.com/ethereum/go-ethereum/common"

func CommonHashToUint32Slice(h common.Hash) []uint32 {
	b := Bytes32{}
	b.FromBytes(h[:])

	return b[:]
}
