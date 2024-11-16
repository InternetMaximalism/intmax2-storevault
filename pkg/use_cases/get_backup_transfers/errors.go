package get_backup_transfers

import "errors"

// ErrUCGetBackupTransfersInputEmpty error: ucGetBackupTransfersInput must not be empty.
var ErrUCGetBackupTransfersInputEmpty = errors.New("ucGetBackupTransfersInput must not be empty")

// ErrGetBackupTransfersFail error: failed to get the backup transfers list.
var ErrGetBackupTransfersFail = errors.New("failed to get the backup transfers list")
