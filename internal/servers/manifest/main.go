package manifest

import (
	"context"
	"net"

	"google.golang.org/grpc"

	pbmanifest "github.com/walkergriggs/enoki/internal/proto/manifest"
	"github.com/walkergriggs/enoki/internal/services/manifest"
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

	return s.Serve(l)
}
