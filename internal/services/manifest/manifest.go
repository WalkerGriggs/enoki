package manifest

import (
	"context"

	pbmanifest "github.com/walkergriggs/enoki/internal/proto/golang/manifest"
)

type ManifestService struct{}

func (s *ManifestService) GetManifest(ctx context.Context, in *pbmanifest.ManifestRequest) (*pbmanifest.Manifest, error) {
	return &pbmanifest.Manifest{
		Id:   1234,
		Path: "/here",
	}, nil
}
