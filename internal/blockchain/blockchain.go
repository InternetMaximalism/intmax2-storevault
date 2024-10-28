package blockchain

import (
	"context"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
)

type serviceBlockchain struct {
	ctx context.Context
	cfg *configs.Config
	log logger.Logger
}

func New(
	ctx context.Context,
	cfg *configs.Config,
	log logger.Logger,
) ServiceBlockchain {
	return &serviceBlockchain{
		ctx: ctx,
		cfg: cfg,
		log: log,
	}
}
