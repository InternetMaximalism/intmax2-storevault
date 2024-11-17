package get_backup_transfer_by_uuid

import "github.com/prodadidb/go-validation"

func (input *UCGetBackupTransferByUuidInput) Valid() error {
	return validation.ValidateStruct(input,
		validation.Field(&input.Recipient, validation.Required),
		validation.Field(&input.Uuid, validation.Required),
	)
}
