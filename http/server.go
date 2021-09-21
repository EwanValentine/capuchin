package http

import (
	"context"
	"net/http"

	"github.com/EwanValentine/capuchin/conf"
	gw "github.com/EwanValentine/capuchin/gen/go/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// NewServer -
func NewServer(conf *conf.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterCapuchinQueryServiceHandlerFromEndpoint(ctx, mux, conf.GRPCAddr, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(conf.HostAddr, mux)
}
