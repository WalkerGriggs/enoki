package storage

import (
	_ "github.com/mattn/go-sqlite3"

	"context"
	"database/sql"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/walkergriggs/enoki/internal/services/storage/db/sqlite"
)

type StorageService struct {
	Queries *sqlite.Queries
}

const (
	srcStoragePath      = "/home/wgriggs/Documents/enoki-assets"
	sqlSchemaPath       = "/home/wgriggs/go/src/github.com/walkergriggs/enoki/internal/services/storage/db/sqlite/schema.sql"
	sqlDriverName       = "sqlite3"
	sqlDriverSourceName = ":memory:"
)

func NewStorageService(ctx context.Context) (*StorageService, error) {
	storagePath, err := filepath.Abs(srcStoragePath)
	if err != nil {
		return nil, err
	}

	ddl, err := ioutil.ReadFile(sqlSchemaPath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(sqlDriverName, sqlDriverSourceName)
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, string(ddl)); err != nil {
		return nil, err
	}

	service := &StorageService{
		Queries: sqlite.New(db),
	}

	if err := service.init(ctx, storagePath); err != nil {
		return nil, err
	}

	return service, nil
}

// TODO(walker): Replace with something robust. This is a hack to get us up and
// running.
func (s *StorageService) init(ctx context.Context, root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) == ".m3u8" {
			base := filepath.Base(path)
			name := base[:len(base)-len(filepath.Ext(path))]

			s.Queries.CreateAsset(ctx, sqlite.CreateAssetParams{
				Path: path,
				Name: name,
			})
		}
		return nil
	})
}

