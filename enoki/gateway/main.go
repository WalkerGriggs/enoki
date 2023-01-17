package gateway

import (
	"context"
	"net/http"
	"log"

	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Println("Dialing")
	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	log.Println("Handling")
	gw, err := newGateway(ctx, conn, []runtime.ServeMuxOption{})
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzServer(conn))
	mux.Handle("/", gw)

	s := &http.Server{
		Addr:    "localhost:8001",
		Handler: allowCORS(mux),
	}

	go func() {
		<-ctx.Done()
		s.Shutdown(context.Background())
	}()

	log.Println("Serving")
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
