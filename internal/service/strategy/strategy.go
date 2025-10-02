package strategy

import (
	"context"

	"github.com/mmedinam1600/product-comparison-api/internal/domain"
)

// Interface define the contract for comparison strategies
type Interface interface {
	// Name returns the identifier name of the strategy
	Name() string

	// ResolveFields determines which fields can be compared according to the strategy.
	// items: products to compare
	// requested: fields requested by the client (nil if not specified)
	// Returns: list of resolved fields that can be compared
	ResolveFields(items []domain.Item, requested *[]string) []string

	// ComputeDiff calculates the differences for each resolved field.
	// items: products to compare
	// resolved: fields to compare (result of ResolveFields)
	// Returns: map of field → differences
	ComputeDiff(ctx context.Context, items []domain.Item, resolved []string) (map[string]domain.DiffField, error)
}
