package gateway

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/walkergriggs/enoki/internal/shared/logging"

	pbmanifest "github.com/walkergriggs/enoki/internal/proto/golang/manifest"
	pbstorage "github.com/walkergriggs/enoki/internal/proto/golang/storage"
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

	logging.WithContext(ctx).Info("Registering manifest service",
		zap.String("Addr", "localhost:8080"),
	)
	err := pbmanifest.RegisterManifestServiceHandlerFromEndpoint(ctx, mux, "localhost:8080", s.DialOpts)
	if err != nil {
		return err
	}

	logging.WithContext(ctx).Info("Registering storage service",
		zap.String("Addr", "localhost:8082"),
	)
	err = pbstorage.RegisterStorageServiceHandlerFromEndpoint(ctx, mux, "localhost:8082", s.DialOpts)
	if err != nil {
		return err
	}

	s.Mux.Handle("/", mux)
	return nil
}

func (s *GatewayService) RegisterHealthzHandler(ctx context.Context) error {
	logging.WithContext(ctx).Info("Registering healthz handler")

	if s.Mux == nil {
		s.Mux = http.NewServeMux()
	}

	s.Mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		logging.WithContext(ctx).Info("Healthz")
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

	logging.WithContext(ctx).Info("Serving gateway http server",
		zap.String("Addr", server.Addr),
	)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
