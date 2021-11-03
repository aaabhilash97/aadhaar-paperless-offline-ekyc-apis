package rest

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	api "github.com/aaabhilash97/aadhaar-paperless-offline-ekyc-apis/pkg/api/v1"
	"github.com/aaabhilash97/aadhaar-paperless-offline-ekyc-apis/pkg/protocol/rest/middleware"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func headerMatcher(headerName string) (mdName string, ok bool) {
	allowedHeaders := [5]string{"appId", "appKey"}
	for _, header := range allowedHeaders {
		if header == headerName {
			return header, true
		}
	}
	return "", false
}

// RunServer runs gRPC service to publish ToDo service
type Opt struct {
	ShutdownCtx context.Context
	GrpcPort    string
	HttpPort    string
	Host        string
	AccessLog   *zap.Logger
	Logger      *zap.Logger
}

func httpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set http status code
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
		w.WriteHeader(code)
	}

	return nil
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}

// allowCORS allows Cross Origin Resource Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// RunServer runs HTTP/REST gateway
func RunServer(opt Opt) error {

	logger := opt.Logger

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(headerMatcher),
		runtime.WithForwardResponseOption(httpResponseModifier),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
				// EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	{
		opts := []grpc.DialOption{grpc.WithInsecure()}
		ctx := context.Background()
		if err := api.RegisterAadhaarServiceHandlerFromEndpoint(ctx, mux, "localhost:"+opt.GrpcPort, opts); err != nil {
			logger.Error("failed to start HTTP gateway: %v", zap.String("reason", err.Error()))
			return err
		}

	}

	srv := &http.Server{
		WriteTimeout:      time.Second * 120,
		ReadHeaderTimeout: time.Second * 120,
		IdleTimeout:       time.Second * 120,
		Addr:              opt.Host + ":" + opt.HttpPort,
		Handler: allowCORS(
			middleware.AddRequestID(
				middleware.AddLogger(opt.AccessLog, mux),
			)),
	}

	go func(shutdownCtx context.Context) {
		<-shutdownCtx.Done()
		logger.Info("shutting down Rest server")
		{
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			err := srv.Shutdown(ctx)
			if err != nil {
				logger.Error("rest.Shutdown", zap.Error(err))
			}
		}

	}(opt.ShutdownCtx)

	logger.Info("starting HTTP/REST gateway", zap.String("port", opt.HttpPort))

	return srv.ListenAndServe()
}
