package get_backup_transaction_by_uuid

import "errors"

// ErrUCGetBackupTransactionByUuidInputEmpty error: ucGetBackupTransactionByUuidInput must not be empty.
var ErrUCGetBackupTransactionByUuidInputEmpty = errors.New("ucGetBackupTransactionByUuidInput must not be empty")

// ErrGetBackupTransactionByIDFail error: failed to get backup transaction by id.
var ErrGetBackupTransactionByIDFail = errors.New("failed to get backup transaction by id")
