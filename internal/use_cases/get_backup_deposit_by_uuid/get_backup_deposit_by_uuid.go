package get_backup_deposit_by_uuid

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_get_backup_deposit_by_uuid.go -package=mocks -source=get_backup_deposit_by_uuid.go

const (
	NotFoundMessage = "Deposit not found."
)

type UCGetBackupDepositByUuid struct {
	Uuid             string
	Recipient        string
	BlockNumber      uint64
	EncryptedDeposit string
	CreatedAt        time.Time
}

type UCGetBackupDepositByUuidInput struct {
	Uuid string `json:"uuid"`
}

type UseCaseGetBackupDepositByUuid interface {
	Do(
		ctx context.Context, input *UCGetBackupDepositByUuidInput,
	) (*UCGetBackupDepositByUuid, error)
}
