package store_vault_server

import (
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"

	"github.com/dimiro1/health"
)

type StoreVaultServer struct {
	log              logger.Logger
	config           *configs.Config
	dbApp            SQLDriverApp
	commands         Commands
	sb               ServiceBlockchain
	cookieForAuthUse bool
	hc               *health.Handler
	vdcs             VerifyDepositConfirmationService
}

// New initializes a new Server struct.
func New(
	log logger.Logger,
	config *configs.Config,
	dbApp SQLDriverApp,
	commands Commands,
	sb ServiceBlockchain,
	cookieForAuthUse bool,
	hc *health.Handler,
	vdcs VerifyDepositConfirmationService,
) *StoreVaultServer {
	const (
		srv  = "storeVaultServer"
		name = "storevaultserver"
	)

	return &StoreVaultServer{
		log:              log.WithFields(logger.Fields{srv: name}),
		config:           config,
		dbApp:            dbApp,
		commands:         commands,
		sb:               sb,
		cookieForAuthUse: cookieForAuthUse,
		hc:               hc,
		vdcs:             vdcs,
	}
}

func (s *StoreVaultServer) Commands() Commands {
	return s.commands
}
