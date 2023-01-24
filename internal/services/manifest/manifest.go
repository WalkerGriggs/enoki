package manifest

import (
	"context"

	"github.com/walkergriggs/enoki/internal/shared/logging"

	pbmanifest "github.com/walkergriggs/enoki/internal/proto/golang/manifest"
)

type ManifestService struct{}

func (s *ManifestService) GetManifest(ctx context.Context, in *pbmanifest.ManifestRequest) (*pbmanifest.Manifest, error) {
	logging.WithContext(ctx).Info("GetManifest")

	return &pbmanifest.Manifest{
		Id:   1234,
		Path: "/here",
	}, nil
}
