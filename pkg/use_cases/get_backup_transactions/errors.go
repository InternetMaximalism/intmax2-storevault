package get_backup_transactions

import "errors"

// ErrUCGetBackupTransactionsInputEmpty error: ucGetBackupTransactionsInput must not be empty.
var ErrUCGetBackupTransactionsInputEmpty = errors.New("ucGetBackupTransactionsInput must not be empty")

// ErrGetBackupTransactionsFail error: failed to get the backup transactions list.
var ErrGetBackupTransactionsFail = errors.New("failed to get the backup transactions list")
