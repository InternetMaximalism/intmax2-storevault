package get_backup_balances

import "errors"

// ErrUCGetBackupBalancesInputEmpty error: ucGetBackupBalancesInput must not be empty.
var ErrUCGetBackupBalancesInputEmpty = errors.New("ucGetBackupBalancesInput must not be empty")

// ErrGetLatestBackupBalanceByUserAddressFail error: failed to get latest the backup balance by user address.
var ErrGetLatestBackupBalanceByUserAddressFail = errors.New(
	"failed to get latest the backup balance by user address",
)
