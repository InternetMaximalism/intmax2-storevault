package get_backup_balances

import "github.com/prodadidb/go-validation"

func (input *UCGetBackupBalancesInput) Valid() error {
	return validation.ValidateStruct(input,
		validation.Field(&input.Sender, validation.Required),
	)
}
