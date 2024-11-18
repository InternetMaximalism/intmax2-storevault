package get_backup_deposit_by_uuid

import "github.com/prodadidb/go-validation"

func (input *UCGetBackupDepositByUuidInput) Valid() error {
	return validation.ValidateStruct(input,
		validation.Field(&input.Uuid, validation.Required),
	)
}
