package middleware

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// codeToLevel redirects OK to DEBUG level logging instead of INFO
// This is example how you can log several gRPC code results
func codeToLevel(code codes.Code) zapcore.Level {
	if code == codes.OK {
		// It is DEBUG
		return zap.DebugLevel
	}
	return grpc_zap.DefaultCodeToLevel(code)
}

// constructLoggingInterceptors return unary and stream interceptors
func constructLoggingInterceptors(logger *zap.Logger) ([]grpc.UnaryServerInterceptor, []grpc.StreamServerInterceptor) {
	// Shared options for the logger, with a custom gRPC code to log level function.
	o := []grpc_zap.Option{
		grpc_zap.WithLevels(codeToLevel),
	}

	// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	// Add unary interceptor
	unaryInterceptors := make([]grpc.UnaryServerInterceptor, 0)
	unaryInterceptors = append(unaryInterceptors,
		grpc_ctxtags.UnaryServerInterceptor(
			grpc_ctxtags.WithFieldExtractor(
				grpc_ctxtags.CodeGenRequestFieldExtractor,
			),
		),
		grpc_zap.UnaryServerInterceptor(logger, o...),
	)

	streamInterceptors := make([]grpc.StreamServerInterceptor, 0)
	streamInterceptors = append(streamInterceptors,
		grpc_ctxtags.StreamServerInterceptor(
			grpc_ctxtags.WithFieldExtractor(
				grpc_ctxtags.CodeGenRequestFieldExtractor,
			),
		),
		grpc_zap.StreamServerInterceptor(logger, o...),
	)

	return unaryInterceptors, streamInterceptors
}
