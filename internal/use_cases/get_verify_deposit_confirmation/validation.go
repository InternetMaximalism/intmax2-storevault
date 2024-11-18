package get_verify_deposit_confirmation

import (
	"errors"

	"github.com/holiman/uint256"
	"github.com/prodadidb/go-validation"
)

// ErrValueInvalid error: value must be valid.
var ErrValueInvalid = errors.New("value must be valid")

func (input *UCGetVerifyDepositConfirmationInput) Valid() error {
	return validation.ValidateStruct(input,
		validation.Field(&input.DepositID, validation.Required, input.IsDepositID()),
	)
}

func (input *UCGetVerifyDepositConfirmationInput) IsDepositID() validation.Rule {
	return validation.By(func(value interface{}) (err error) {
		v, ok := value.(string)
		if !ok {
			return ErrValueInvalid
		}

		var depositID uint256.Int
		err = depositID.Scan(v)
		if err != nil {
			return ErrValueInvalid
		}

		input.ConvertDepositID = &depositID

		return nil
	})
}
