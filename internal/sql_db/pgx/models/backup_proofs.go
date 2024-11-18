package models

import (
	"time"
)

type BackupSenderProof struct {
	ID                         string
	EnoughBalanceProofBodyHash string
	LastBalanceProofBody       []byte
	BalanceTransitionProofBody []byte
	CreatedAt                  time.Time
}
