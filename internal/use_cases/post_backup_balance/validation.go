package post_backup_balance

import "github.com/prodadidb/go-validation"

func (input *UCPostBackupBalanceInput) Valid() error {
	return validation.ValidateStruct(input,
		validation.Field(&input.User, validation.Required),
		validation.Field(&input.EncryptedBalanceProof, validation.Required),
		validation.Field(&input.EncryptedBalanceData, validation.Required),
		// validation.Field(&input.EncryptedTxs, validation.Required),
		// validation.Field(&input.EncryptedTransfers, validation.Required),
		// validation.Field(&input.EncryptedDeposits, validation.Required),
		validation.Field(&input.Signature, validation.Required),
	)
}
