package types

import (
	intMaxTree "intmax2-store-vault/internal/tree"
	intMaxTypes "intmax2-store-vault/internal/types"
)

const (
	int1Key  = 1
	int2Key  = 2
	int3Key  = 3
	int4Key  = 4
	int8Key  = 8
	int32Key = 32
	int47Key = 47

	numTransfersInTx       = int1Key << intMaxTree.TRANSFER_TREE_HEIGHT
	balancePublicInputsLen = int47Key
)

type poseidonHashOut = intMaxTypes.PoseidonHashOut
