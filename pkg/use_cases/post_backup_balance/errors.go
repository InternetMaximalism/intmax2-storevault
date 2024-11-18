package post_backup_balance

import "errors"

// ErrUCPostBackupBalanceInputEmpty error: ucPostBackupBalanceInput must not be empty.
var ErrUCPostBackupBalanceInputEmpty = errors.New("ucPostBackupBalanceInput must not be empty")

// ErrCreateBackupBalanceFail error: failed to create the backup balance.
var ErrCreateBackupBalanceFail = errors.New("failed to create the backup balance")
