package post_backup_transfer

import "errors"

// ErrUCPostBackupTransferInputEmpty error: ucPostBackupTransferInput must not be empty.
var ErrUCPostBackupTransferInputEmpty = errors.New("ucPostBackupTransferInput must not be empty")

// ErrCreateBackupTransferWithDBFail error: failed to create new backup transfer with DB.
var ErrCreateBackupTransferWithDBFail = errors.New(
	"failed to create new backup transfer with DB",
)
