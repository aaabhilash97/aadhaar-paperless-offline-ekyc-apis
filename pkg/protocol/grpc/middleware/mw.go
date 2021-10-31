package middleware

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type MWOptions struct {
	AccessLog *zap.Logger
}

func AddMiddlewares(opt MWOptions, grpcOptions []grpc.ServerOption) []grpc.ServerOption {
	unaryInterceptors := make([]grpc.UnaryServerInterceptor, 0)
	streamInterceptors := make([]grpc.StreamServerInterceptor, 0)

	// Construct logging interceptor
	unaryLog, streamLog := constructLoggingInterceptors(opt.AccessLog)
	unaryInterceptors = append(unaryInterceptors, unaryLog...)
	streamInterceptors = append(streamInterceptors, streamLog...)

	// // // Construct sentry interceptor
	// unarySentry := constructSentryInterceptors()
	// unaryInterceptors = append(unaryInterceptors, unarySentry...)

	// Add unary interceptors
	grpcOptions = append(grpcOptions, grpc_middleware.WithUnaryServerChain(
		unaryInterceptors...,
	))

	// Add stream interceptors
	grpcOptions = append(grpcOptions, grpc_middleware.WithStreamServerChain(
		streamInterceptors...,
	))

	return grpcOptions
}
