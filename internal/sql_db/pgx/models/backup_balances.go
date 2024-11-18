package models

import "time"

type BackupBalance struct {
	ID                    string
	UserAddress           string
	EncryptedBalanceProof string
	EncryptedBalanceData  string
	EncryptedTxs          []string
	EncryptedTransfers    []string
	EncryptedDeposits     []string
	Signature             string
	BlockNumber           uint64
	CreatedAt             time.Time
}
