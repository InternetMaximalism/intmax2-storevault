package main

import (
	"context"
	"intmax2-store-vault/cmd/migrator"
	"intmax2-store-vault/cmd/store_vault_server"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/blockchain"
	"intmax2-store-vault/internal/cli"
	"intmax2-store-vault/internal/open_telemetry"
	"intmax2-store-vault/pkg/logger"
	"intmax2-store-vault/pkg/sql_db/db_app"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dimiro1/health"
)

func main() {
	cfg := configs.New()
	log := logger.New(cfg.LOG.Level, cfg.LOG.TimeFormat, cfg.LOG.JSON, cfg.LOG.IsLogLine)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	const int1 = 1
	done := make(chan os.Signal, int1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer close(done)

	go func() {
		<-done
		const msg = "SIGTERM detected"
		log.Errorf(msg)
		if cancel != nil {
			cancel()
		}
	}()

	err := open_telemetry.Init(cfg.OpenTelemetry.Enable)
	if err != nil {
		const msg = "open_telemetry init: %v"
		log.Errorf(msg, err)
		return
	}

	var dbApp db_app.SQLDb
	dbApp, err = db_app.New(ctx, log, &cfg.SQLDb)
	if err != nil {
		const msg = "db application init: %v"
		log.Errorf(msg, err)
		return
	}

	hc := health.NewHandler()
	bc := blockchain.New(ctx, cfg, log)

	wg := sync.WaitGroup{}

	err = cli.Run(
		ctx,
		migrator.NewMigratorCmd(ctx, log, dbApp),
		store_vault_server.NewServerCmd(&store_vault_server.StoreVaultServer{
			Context: ctx,
			Cancel:  cancel,
			Config:  cfg,
			Log:     log,
			SB:      bc,
			DbApp:   dbApp,
			WG:      &wg,
			HC:      &hc,
		}),
	)
	if err != nil {
		const msg = "cli: %v"
		log.Errorf(msg, err)
		return
	}

	wg.Wait()
}
