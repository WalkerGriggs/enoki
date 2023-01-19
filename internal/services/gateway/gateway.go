package gateway

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/walkergriggs/enoki/internal/proto/manifests/v1"
)

type GatewayService struct {
	Mux       *http.ServeMux
	Server    *http.Server
	ServeOpts []runtime.ServeMuxOption
	DialOpts  []grpc.DialOption
}

func (s *GatewayService) RegisterGatewayHandlers(ctx context.Context) error {
	if s.Mux == nil {
		s.Mux = http.NewServeMux()
	}

	mux := runtime.NewServeMux(s.ServeOpts...)
	err := v1.RegisterManifestServiceHandlerFromEndpoint(ctx, mux, "localhost:8080", s.DialOpts)
	if err != nil {
		return err
	}

	s.Mux.Handle("/", mux)
	return nil
}

func (s *GatewayService) RegisterHealthzHandler(ctx context.Context) error {
	if s.Mux == nil {
		s.Mux = http.NewServeMux()
	}

	s.Mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("ok"))
	})
	return nil
}

func ListenAndServe(ctx context.Context, service *GatewayService) error {
	allowCORS := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)

			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				headers := []string{"Content-Type", "Accept", "Authorization"}
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
				
				methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
				return
			}
		}
		service.Mux.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:    "localhost:8001",
		Handler: allowCORS,
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
