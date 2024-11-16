package post_backup_transaction

import "errors"

// ErrUCPostBackupTransactionInputEmpty error: ucPostBackupTransactionInput must not be empty.
var ErrUCPostBackupTransactionInputEmpty = errors.New("ucPostBackupTransactionInput must not be empty")

// ErrCreateBackupTransactionFail error: failed to create the backup transaction.
var ErrCreateBackupTransactionFail = errors.New("failed to create the backup transaction")
