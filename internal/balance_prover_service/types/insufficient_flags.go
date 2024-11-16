package types

import (
	"encoding/binary"

	"github.com/iden3/go-iden3-crypto/ffg"
)

const (
	InsufficientFlagsLen = numTransfersInTx / int32Key
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
		if y >= uint64(int1Key)<<int32Key {
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
		flags.Limbs[limbIndex] |= int1Key << bitIndex
	} else {
		flags.Limbs[limbIndex] &^= int1Key << bitIndex
	}
}

func (flags *InsufficientFlags) RandomAccess(index int) bool {
	limbIndex := index / int32Key
	bitIndex := index % int32Key

	return flags.Limbs[limbIndex]&(int1Key<<bitIndex) != 0
}

func (flags *InsufficientFlags) Bytes() []byte {
	buf := make([]byte, InsufficientFlagsLen*int4Key)
	for i, limb := range flags.Limbs {
		binary.BigEndian.PutUint32(buf[i*int4Key:], limb)
	}

	return buf
}
