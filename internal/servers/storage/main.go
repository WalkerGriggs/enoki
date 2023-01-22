package storage

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/walkergriggs/enoki/internal/services/storage"

	pbstorage "github.com/walkergriggs/enoki/internal/proto/golang/storage"
)

func Run(ctx context.Context) error {
	l, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		return err
	}

	defer func() {
		l.Close()
	}()

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

	service, err := storage.NewStorageService(ctx)
	if err != nil {
		return err
	}

	pbstorage.RegisterStorageServiceServer(s, service)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	return s.Serve(l)
}
