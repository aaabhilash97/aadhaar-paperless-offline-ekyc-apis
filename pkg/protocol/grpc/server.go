package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/aaabhilash97/aadhaar_scrapper_apis/pkg/protocol/grpc/middleware"
	service "github.com/aaabhilash97/aadhaar_scrapper_apis/pkg/service/v1"

	api "github.com/aaabhilash97/aadhaar_scrapper_apis/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerOpt interface {
	GetLogger() *zap.Logger
	GetAccessLogger() *zap.Logger
}

// RunServer runs gRPC service to publish ToDo service
type Opt struct {
	ShutdownCtx       context.Context
	AadhaarServiceOpt service.NewAadhaarServiceI
	ServerOpt         ServerOpt
	Port              string
	Host              string
}

func RunServer(opt Opt) error {
	logger := opt.ServerOpt.GetLogger()

	listen, err := net.Listen("tcp", opt.Host+":"+opt.Port)
	if err != nil {
		return err
	}

	serverOptions := []grpc.ServerOption{}

	// add middleware
	serverOptions = middleware.AddMiddlewares(middleware.MWOptions{
		AccessLog: opt.ServerOpt.GetAccessLogger(),
	}, serverOptions)
	// register service

	serverOptions = append(serverOptions, grpc.MaxRecvMsgSize(1024*1024*10))
	server := grpc.NewServer(serverOptions...)

	API, err := service.NewAadhaarService(opt.AadhaarServiceOpt)
	if err != nil {
		logger.Error("NewService", zap.Error(err))
		return err
	}
	api.RegisterAadhaarServiceServer(server, API)

	// API2 := service.NewService2(opt.Conf)
	// api.RegisterAadhaarService2Server(server, API2)

	reflection.Register(server)

	// graceful shutdown

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, os.Interrupt)
	go func(shutdownCtx context.Context) {
		<-shutdownCtx.Done()
		logger.Info("shutting down gRPC server")
		server.GracefulStop()
	}(opt.ShutdownCtx)

	// start gRPC server
	logger.Info("starting gRPC server", zap.String("port", opt.Port))
	return server.Serve(listen)
}
