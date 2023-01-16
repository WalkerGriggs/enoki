package manifest

import (
	"context"

	"github.com/walkergriggs/enoki/proto/manifests/v1"
)

type ManifestService struct{}

func (s *ManifestService) GetManifest(ctx context.Context, in *v1.ManifestRequest) (*v1.Manifest, error) {
	return &v1.Manifest{
		Id:   1234,
		Path: "/here",
	}, nil
}
