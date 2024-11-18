package get_backup_transfer_by_uuid

import "errors"

// ErrUCGetBackupTransferByUuidInputEmpty error: ucGetBackupTransferByUuidInput must not be empty.
var ErrUCGetBackupTransferByUuidInputEmpty = errors.New("ucGetBackupTransferByUuidInput must not be empty")

// ErrGetBackupTransferByIDFail error: failed to get backup transfer by id.
var ErrGetBackupTransferByIDFail = errors.New("failed to get backup transfer by id")
