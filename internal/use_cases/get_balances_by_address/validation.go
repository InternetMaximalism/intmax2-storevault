package get_balances_by_address

import "github.com/prodadidb/go-validation"

func (input *UCGetBalancesByAddressInput) Valid() error {
	return validation.ValidateStruct(input,
		validation.Field(&input.Address, validation.Required),
	)
}
