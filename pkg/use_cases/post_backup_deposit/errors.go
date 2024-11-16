package post_backup_deposit

import "errors"

// ErrUCPostBackupDepositInputEmpty error: ucPostBackupDepositInput must not be empty.
var ErrUCPostBackupDepositInputEmpty = errors.New("ucPostBackupDepositInput must not be empty")

// ErrCreateBackupDepositFail error: failed to create the backup deposit.
var ErrCreateBackupDepositFail = errors.New("failed to create the backup deposit")
