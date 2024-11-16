package get_backup_deposit_by_hash

import "errors"

// ErrUCGetBackupDepositByHashInputEmpty error: ucGetBackupDepositByHashInput must not be empty.
var ErrUCGetBackupDepositByHashInputEmpty = errors.New("ucGetBackupDepositByHashInput must not be empty")

// ErrGetBackupDepositByRecipientAndDepositDoubleHash error: failed to get backup deposit by recipient and deposit double hash.
var ErrGetBackupDepositByRecipientAndDepositDoubleHash = errors.New(
	"failed to get backup deposit by recipient and deposit double hash",
)
