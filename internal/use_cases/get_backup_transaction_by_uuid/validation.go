package get_backup_transaction_by_uuid

import "github.com/prodadidb/go-validation"

func (input *UCGetBackupTransactionByUuidInput) Valid() error {
	return validation.ValidateStruct(input,
		validation.Field(&input.Sender, validation.Required),
		validation.Field(&input.Uuid, validation.Required),
	)
}
