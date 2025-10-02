package data

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/goccy/go-json"
	"github.com/mmedinam1600/product-comparison-api/internal/domain"
	"go.uber.org/zap"
)

// CatalogRepository define the contract to access the catalog of products
type CatalogRepository interface {
	// GetByIDs retrieves items by their IDs.
	// Returns: items found, IDs that were not found
	GetByIDs(ctx context.Context, ids []string) ([]domain.Item, []string)

	// GetAll retrieves all items from the catalog
	GetAll(ctx context.Context) []domain.Item
}

// FileCatalogRepo implements CatalogRepository loading data from a JSON file
type FileCatalogRepo struct {
	logger   *zap.Logger
	filePath string
	// atomic.Value for lock-free catalog reads
	catalog atomic.Value // map[string]domain.Item
}

func NewFileCatalogRepo(filePath string, logger *zap.Logger) (*FileCatalogRepo, error) {
	repo := &FileCatalogRepo{
		logger:   logger,
		filePath: filePath,
	}

	if err := repo.reload(); err != nil {
		return nil, fmt.Errorf("failed to load catalog: %w", err)
	}

	logger.Info("catalog loaded successfully",
		zap.String("file", filePath),
		zap.Int("items", len(repo.getCatalog())),
	)

	return repo, nil
}

// GetByIDs implementa CatalogRepository.GetByIDs
func (r *FileCatalogRepo) GetByIDs(ctx context.Context, ids []string) ([]domain.Item, []string) {
	catalog := r.getCatalog()

	found := make([]domain.Item, 0, len(ids))
	missing := make([]string, 0)

	// Maintain the order of the requested IDs
	for _, id := range ids {
		if item, exists := catalog[id]; exists {
			found = append(found, item)
		} else {
			missing = append(missing, id)
		}
	}

	return found, missing
}

// GetAll implements CatalogRepository.GetAll
func (r *FileCatalogRepo) GetAll(ctx context.Context) []domain.Item {
	catalog := r.getCatalog()

	items := make([]domain.Item, 0, len(catalog))
	for _, item := range catalog {
		items = append(items, item)
	}

	return items
}

// reload loads (or reloads) the catalog from the JSON file
func (r *FileCatalogRepo) reload() error {
	// 1. Read file
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// 2. Parse JSON to array of items
	var items []domain.Item
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// 3. Build index by ID
	catalog := make(map[string]domain.Item, len(items))
	for _, item := range items {
		catalog[item.ID] = item
	}

	// 4. Update catalog atomically
	r.catalog.Store(catalog)

	return nil
}

// getCatalog gets the current catalog in a thread-safe way
func (r *FileCatalogRepo) getCatalog() map[string]domain.Item {
	return r.catalog.Load().(map[string]domain.Item)
}
