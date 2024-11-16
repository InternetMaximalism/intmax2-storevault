package get_balances_by_address

import "errors"

// ErrUCGetBalancesByAddressInputEmpty error: ucGetBalancesByAddressInput must not be empty.
var ErrUCGetBalancesByAddressInputEmpty = errors.New("ucGetBalancesByAddressInput must not be empty")

// ErrGetBackupDepositsFail error: failed to ge the backup deposits list.
var ErrGetBackupDepositsFail = errors.New("failed to ge the backup deposits list")

// ErrGetBackupTransactionsFail error: failed to get the backup transactions list.
var ErrGetBackupTransactionsFail = errors.New("failed to get the backup transactions list")

// ErrGetBackupTransfersFail error: failed to get the backup transfers list.
var ErrGetBackupTransfersFail = errors.New("failed to get the backup transfers list")
