package get_backup_transaction_by_hash

import "errors"

// ErrUCGetBackupTransactionByHashInputEmpty error: ucGetBackupTransactionByHashInput must not be empty.
var ErrUCGetBackupTransactionByHashInputEmpty = errors.New("ucGetBackupTransactionByHashInput must not be empty")

// ErrGetBackupTransactionBySenderAndTxDoubleHashFail error: failed to get backup transaction by sender and tx double hash.
var ErrGetBackupTransactionBySenderAndTxDoubleHashFail = errors.New(
	"failed to get backup transaction by sender and tx double hash",
)
