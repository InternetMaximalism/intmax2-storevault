package store_vault_server

import (
	"context"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/configs/buildvars"
	"intmax2-store-vault/docs/swagger"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/pb/gateway"
	"intmax2-store-vault/internal/pb/gateway/consts"
	"intmax2-store-vault/internal/pb/gateway/http_response_modifier"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	"intmax2-store-vault/internal/pb/listener"
	"intmax2-store-vault/internal/verify_deposit_confirmation_service"
	server "intmax2-store-vault/pkg/grpc_server/store_vault_server"
	"intmax2-store-vault/third_party"
	"sync"
	"time"

	"github.com/dimiro1/health"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type StoreVaultServer struct {
	Context context.Context
	Cancel  context.CancelFunc
	WG      *sync.WaitGroup
	Config  *configs.Config
	Log     logger.Logger
	DbApp   SQLDriverApp
	SB      ServiceBlockchain
	HC      *health.Handler
}

func NewServerCmd(s *StoreVaultServer) *cobra.Command {
	const (
		use   = "run"
		short = "Run store vault server"
	)
	return &cobra.Command{
		Use:   use,
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			err := s.SB.SetupEthereumNetworkChainID(s.Context)
			if err != nil {
				const msg = "failed to setup ethereum network chain ID: %+v"
				s.Log.Fatalf(msg, err.Error())
			}

			err = s.SB.SetupScrollNetworkChainID(s.Context)
			if err != nil {
				const msg = "failed to setup scroll network chain ID: %+v"
				s.Log.Fatalf(msg, err.Error())
			}

			var vdcs verify_deposit_confirmation_service.VerifyDepositConfirmationService
			vdcs, err = verify_deposit_confirmation_service.New(s.Context, s.Config, s.Log, s.SB)
			if err != nil {
				const msg = "the verify-deposit-confirmation service init: %+v"
				s.Log.Fatalf(msg, err.Error())
			}

			if err = s.Init(vdcs); err != nil {
				const msg = "failed to start api: %+v"
				s.Log.Fatalf(msg, err.Error())
			}
		},
	}
}

// TODO: Common parts should be shared between the server side and the client side.
func (s *StoreVaultServer) Init(
	vdcs VerifyDepositConfirmationService,
) error {
	tm := time.Duration(s.Config.HTTP.Timeout) * time.Second

	var c *cors.Cors
	if s.Config.HTTP.CORSAllowAll {
		c = cors.AllowAll()
	} else {
		c = cors.New(cors.Options{
			AllowedOrigins:       s.Config.HTTP.CORS,
			AllowedMethods:       s.Config.HTTP.CORSAllowMethods,
			AllowedHeaders:       s.Config.HTTP.CORSAllowHeaders,
			ExposedHeaders:       s.Config.HTTP.CORSExposeHeaders,
			AllowCredentials:     s.Config.HTTP.CORSAllowCredentials,
			MaxAge:               s.Config.HTTP.CORSMaxAge,
			OptionsSuccessStatus: s.Config.HTTP.CORSStatusCode,
		})
	}

	srv := server.New(
		s.Log, s.Config, s.DbApp, server.NewCommands(), s.SB, s.Config.HTTP.CookieForAuthUse, s.HC, vdcs,
	)
	ctx := context.WithValue(s.Context, consts.AppConfigs, s.Config)

	const (
		version   = "version"
		buildtime = "buildtime"
		app       = "app"
		appName   = " (node) "
		sqlDBApp  = "sql-db-app"
	)

	// run externals gRPC server listener
	grpcErr, gRPCServerStop := listener.Run(
		ctx,
		s.Log,
		appName,
		s.Config.GRPC.Addr(), // listen incoming host:port for gRPC server
		func(s grpc.ServiceRegistrar) {
			node.RegisterInfoServiceServer(s, srv)
			node.RegisterStoreVaultServiceServer(s, srv)
		},
	)

	// healthCheck
	s.HC.AddChecker(sqlDBApp, s.DbApp)
	s.HC.AddInfo(app, map[string]any{
		version:   buildvars.Version,
		buildtime: buildvars.BuildTime,
	})

	// run web -> gRPC gateway
	gw, grpcGwErr := gateway.Run(
		ctx,
		&gateway.Params{
			Name:               appName,
			Logger:             s.Log,
			GatewayAddr:        s.Config.HTTP.Addr(), // listen incoming host:port for rest api
			DialAddr:           s.Config.GRPC.Addr(), // connect to gRPC server host:port
			HTTPTimeout:        tm,
			HealthCheckHandler: s.HC,
			Services: []gateway.RegisterServiceHandlerFunc{
				node.RegisterInfoServiceHandler,
				node.RegisterStoreVaultServiceHandler,
			},
			CorsHandler: c.Handler,
			Swagger: &gateway.Swagger{
				HostURL:            s.Config.Swagger.HostURL,
				BasePath:           s.Config.Swagger.BasePath,
				SwaggerPath:        configs.SwaggerStoreVaultPath,
				FsSwagger:          swagger.FsSwaggerStoreVault,
				OpenAPIPath:        configs.SwaggerOpenAPIStoreVaultPath,
				FsOpenAPI:          third_party.OpenAPIStoreVault,
				RegexpBuildVersion: s.Config.Swagger.RegexpBuildVersion,
				RegexpHostURL:      s.Config.Swagger.RegexpHostURL,
				RegexpBasePATH:     s.Config.Swagger.RegexpBasePATH,
			},
			Cookies: &http_response_modifier.Cookies{
				ForAuthUse:         s.Config.HTTP.CookieForAuthUse,
				Secure:             s.Config.HTTP.CookieSecure,
				Domain:             s.Config.HTTP.CookieDomain,
				SameSiteStrictMode: s.Config.HTTP.CookieSameSiteStrictMode,
			},
		},
	)

	const (
		start  = "%sapplication started (version: %s buildtime: %s)"
		finish = "%sapplication finished"
	)

	s.Log.Infof(start, appName, buildvars.Version, buildvars.BuildTime)
	defer s.Log.Infof(finish, appName)

	var err error
	select {
	case <-s.Context.Done():
	case err = <-grpcErr:
		const msg = "%sgRPC server error: %s"
		s.Log.Errorf(msg, appName, err)
	case err = <-grpcGwErr:
		const msg = "%sgRPC gateway error: %s"
		s.Log.Errorf(msg, appName, err)
	}

	if gw != nil {
		gw.SetStatus(health.Down)
	}

	gRPCServerStop()
	s.Cancel()

	return nil
}
