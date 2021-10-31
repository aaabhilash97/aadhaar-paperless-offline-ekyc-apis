package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaabhilash97/aadhaar_scrapper_apis/internal/appconfig"
	grpcServer "github.com/aaabhilash97/aadhaar_scrapper_apis/pkg/protocol/grpc"
	"github.com/aaabhilash97/aadhaar_scrapper_apis/pkg/protocol/rest"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// RunServer runs gRPC server and HTTP gateway
func RunServer(conf *appconfig.Config) error {

	if len(conf.Server.GrpcPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", conf.Server.GrpcPort)
	}

	if len(conf.Server.HttpPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", conf.Server.HttpPort)
	}

	shutdownCtx, shutdownCtxCancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, os.Interrupt)
	go func() {
		for range c {
			shutdownCtxCancel()
		}
	}()

	logger := conf.GetLogger()
	accessLogger := conf.GetAccessLogger()

	go func() {
		err := rest.RunServer(rest.Opt{
			ShutdownCtx: shutdownCtx,
			Host:        conf.Server.Host,
			GrpcPort:    conf.Server.GrpcPort,
			HttpPort:    conf.Server.HttpPort,
			AccessLog:   accessLogger,
			Logger:      logger,
		})
		if err != nil {
			if err == http.ErrServerClosed {
				logger.Error("rest.RunServer", zap.NamedError("Closed", err))
			} else {
				logger.Error("rest.RunServer", zap.Error(err))
			}
		}
	}()

	if err := grpcServer.RunServer(grpcServer.Opt{
		ShutdownCtx: shutdownCtx,
		Host:        conf.Server.Host,
		Port:        conf.Server.GrpcPort,

		AadhaarServiceOpt: conf,
		ServerOpt:         conf,
	}); err != nil {
		if err == grpc.ErrServerStopped {
			logger.Error("grpcServ.RunServer", zap.NamedError("Closed", err))
		} else {
			logger.Error("grpcServ.RunServer", zap.Error(err))
		}
		return err
	}
	return nil
}
