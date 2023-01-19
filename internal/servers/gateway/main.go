package gateway

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/walkergriggs/enoki/internal/services/gateway"
)

const (
	ManifestAddress string = "localhost:8080"
	GatewayAddress  string = "localhost:8001"
)

func Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	service := &gateway.GatewayService{
		Mux:       http.NewServeMux(),
		ServeOpts: []runtime.ServeMuxOption{},
		DialOpts:  []grpc.DialOption{grpc.WithInsecure()},
	}

	if err := service.RegisterGatewayHandlers(ctx); err != nil {
		return err
	}

	if err := service.RegisterHealthzHandler(ctx); err != nil {
		return err
	}

	return gateway.ListenAndServe(ctx, service)
}
