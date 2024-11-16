package get_backup_transfer_by_hash

import "errors"

// ErrUCGetBackupTransferByHashInputEmpty error: ucGetBackupTransferByHashInput must not be empty.
var ErrUCGetBackupTransferByHashInputEmpty = errors.New("ucGetBackupTransferByHashInput must not be empty")

// ErrGetBackupTransferByRecipientAndTransferDoubleHashFail error: failed to get backup transfer by recipient and transfer double hash.
var ErrGetBackupTransferByRecipientAndTransferDoubleHashFail = errors.New(
	"failed to get backup transfer by recipient and transfer double hash",
)
