package post_backup_transaction

import (
	"encoding/base64"
	"errors"

	"github.com/prodadidb/go-validation"
)

// ErrValueInvalid error: value must be valid.
var ErrValueInvalid = errors.New("value must be valid")

func (input *UCPostBackupTransactionInput) Valid() error {
	return validation.ValidateStruct(input,
		validation.Field(&input.TxHash, validation.Required),
		validation.Field(&input.EncryptedTx, validation.Required),
		validation.Field(&input.SenderEnoughBalanceProofBody,
			validation.Required,
			validation.By(func(value interface{}) error {
				var isNil bool
				value, isNil = validation.Indirect(value)

				if isNil || validation.IsEmpty(value) {
					return nil
				}

				body, ok := value.(UCPostBackupTransactionInputEnoughBalanceProofBody)
				if !ok {
					return ErrValueInvalid
				}

				return validation.ValidateStruct(&body,
					validation.Field(&body.TransferStepProofBody,
						validation.Required,
						validation.By(func(value interface{}) (err error) {
							vTransferStepProofBody, okTransferStepProofBody := value.(string)
							if !okTransferStepProofBody {
								return ErrValueInvalid
							}

							input.ConvertSenderEnoughBalanceProofBody.ConvertTransferStepProofBody, err =
								base64.StdEncoding.DecodeString(vTransferStepProofBody)
							if err != nil {
								return ErrValueInvalid
							}

							return nil
						}),
					),
					validation.Field(&body.PrevBalanceProofBody,
						validation.Required,
						validation.By(func(value interface{}) (err error) {
							vPrevBalanceProofBody, okPrevBalanceProofBody := value.(string)
							if !okPrevBalanceProofBody {
								return ErrValueInvalid
							}

							input.ConvertSenderEnoughBalanceProofBody.ConvertPrevBalanceProofBody, err =
								base64.StdEncoding.DecodeString(vPrevBalanceProofBody)
							if err != nil {
								return ErrValueInvalid
							}

							return nil
						}),
					),
				)
			}),
		),
		validation.Field(&input.Sender, validation.Required),
		validation.Field(&input.BlockNumber, validation.Required),
		validation.Field(&input.Signature, validation.Required),
	)
}
