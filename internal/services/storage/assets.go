package storage

import (
	"context"
	
	pbasset "github.com/walkergriggs/enoki/internal/proto/golang/asset"
	pbstorage "github.com/walkergriggs/enoki/internal/proto/golang/storage"
)

func (s *StorageService) GetAsset(ctx context.Context, in *pbstorage.AssetRequest) (*pbstorage.AssetResponse, error) {
	asset, err := s.Queries.GetAsset(ctx, in.AssetId)
	if err != nil {
		return nil, err
	}

	return &pbstorage.AssetResponse{
		Asset: &pbasset.Asset{
			Id:   uint64(asset.ID),
			Name: asset.Name,
			Path: asset.Path,
		},
	}, nil
}
