package types

import "github.com/iden3/go-iden3-crypto/ffg"

const (
	int4Key  = 4
	int5Key  = 5
	int8Key  = 8
	int32Key = 32
)

func FieldElementSliceToUint32Slice(value []ffg.Element) []uint32 {
	v := make([]uint32, len(value))
	for i, x := range value {
		y := x.ToUint64Regular()
		if y >= uint64(1)<<int32Key {
			panic("overflow")
		}
		v[i] = uint32(y)
	}

	return v
}
