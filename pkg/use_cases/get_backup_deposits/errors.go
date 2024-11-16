package get_backup_deposits

import "errors"

// ErrUCGetBackupDepositsInputEmpty error: ucGetBackupDepositsInput must not be empty.
var ErrUCGetBackupDepositsInputEmpty = errors.New("ucGetBackupDepositsInput must not be empty")

// ErrGetBackupDepositsFail error: failed to get backup deposits list.
var ErrGetBackupDepositsFail = errors.New("failed to get backup deposits list")
