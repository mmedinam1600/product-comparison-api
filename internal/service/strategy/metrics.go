package strategy

import (
	"strings"

	"github.com/mmedinam1600/product-comparison-api/internal/domain"
)

// fieldMetrics maps field names to their default metrics
var fieldMetrics = map[string]domain.Metric{
	// Root fields where lower is better
	"price": domain.LowerIsBetter,

	// Root fields where higher is better
	"rating": domain.HigherIsBetter,

	// Specifications fields where lower is better
	"specifications.weight": domain.LowerIsBetter,

	// Specifications fields where higher is better
	"specifications.sensor_dpi":   domain.HigherIsBetter,
	"specifications.buttons":      domain.HigherIsBetter,
	"specifications.battery_life": domain.HigherIsBetter,
	"specifications.screen_size":  domain.HigherIsBetter,
	"specifications.refresh_rate": domain.HigherIsBetter,

	// Boolean fields where true is better
	"specifications.wireless":         domain.TrueIsBetter,
	"specifications.noise_cancelling": domain.TrueIsBetter,
	"specifications.backlit":          domain.TrueIsBetter,
}

// GetMetricForField returns the appropriate metric for a given field.
// If the field does not have a predefined metric, returns nil.
func GetMetricForField(fieldPath string) *domain.Metric {
	// Normalize the path (just in case)
	normalizedPath := strings.ToLower(strings.TrimSpace(fieldPath))

	if metric, exists := fieldMetrics[normalizedPath]; exists {
		return &metric
	}

	return nil
}
