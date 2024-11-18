package get_backup_deposit_by_uuid

import "errors"

// ErrUCGetBackupDepositByUuidInputEmpty error: ucGetBackupDepositByUuidInput must not be empty.
var ErrUCGetBackupDepositByUuidInputEmpty = errors.New("ucGetBackupDepositByUuidInput must not be empty")

// ErrGetBackupDepositByIDFail error: failed to get backup deposit by id.
var ErrGetBackupDepositByIDFail = errors.New("failed to get backup deposit by")
