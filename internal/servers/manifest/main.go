package manifest

import (
	"context"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/walkergriggs/enoki/internal/services/manifest"
	"github.com/walkergriggs/enoki/internal/shared/logging"

	pbmanifest "github.com/walkergriggs/enoki/internal/proto/golang/manifest"
)

func Run(ctx context.Context) error {
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}

	defer func() {
		l.Close()
	}()

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

	pbmanifest.RegisterManifestServiceServer(s, &manifest.ManifestService{})

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	logging.WithContext(ctx).Info("Serving manifest server",
		zap.String("Addr", "localhost:8080"),
	)

	return s.Serve(l)
}
