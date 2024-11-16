package post_backup_transaction

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate mockgen -destination=../mocks/mock_post_backup_transaction.go -package=mocks -source=post_backup_transaction.go

const (
	SuccessMsg = "Backup transaction accepted."
)

type UCPostBackupTransactionInputConvertEnoughBalanceProofBody struct {
	ConvertPrevBalanceProofBody  []byte `json:"-"`
	ConvertTransferStepProofBody []byte `json:"-"`
}

func (proof *UCPostBackupTransactionInputConvertEnoughBalanceProofBody) Hash() string {
	var buf []byte
	buf = append(buf, proof.ConvertPrevBalanceProofBody...)
	buf = append(buf, proof.ConvertTransferStepProofBody...)
	output := crypto.Keccak256(buf)

	return hexutil.Encode(output)
}

type UCPostBackupTransactionInputEnoughBalanceProofBody struct {
	PrevBalanceProofBody         string `json:"prevBalanceProof"`
	ConvertPrevBalanceProofBody  []byte `json:"-"`
	TransferStepProofBody        string `json:"transferStepProof"`
	ConvertTransferStepProofBody []byte `json:"-"`
}

type UCPostBackupTransactionInput struct {
	TxHash                              string                                                     `json:"txHash"`
	EncryptedTx                         string                                                     `json:"encryptedTx"`
	SenderEnoughBalanceProofBody        *UCPostBackupTransactionInputEnoughBalanceProofBody        `json:"senderEnoughBalanceProofBody"`
	ConvertSenderEnoughBalanceProofBody *UCPostBackupTransactionInputConvertEnoughBalanceProofBody `json:"-"`
	Sender                              string                                                     `json:"sender"`
	BlockNumber                         uint32                                                     `json:"blockNumber"`
	Signature                           string                                                     `json:"signature"`
}

// UseCasePostBackupTransaction describes PostBackupTransaction contract.
type UseCasePostBackupTransaction interface {
	Do(ctx context.Context, input *UCPostBackupTransactionInput) error
}
