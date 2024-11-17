package get_backup_deposit_by_uuid

import "errors"

// ErrUCGetBackupDepositByUuidInputEmpty error: ucGetBackupDepositByUuidInput must not be empty.
var ErrUCGetBackupDepositByUuidInputEmpty = errors.New("ucGetBackupDepositByUuidInput must not be empty")

// ErrGetBackupDepositByIDAndRecipientFail error: failed to get backup deposit by id and recipient.
var ErrGetBackupDepositByIDAndRecipientFail = errors.New(
	"failed to get backup deposit by id and recipient",
)
