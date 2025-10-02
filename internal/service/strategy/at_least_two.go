package strategy

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/mmedinam1600/product-comparison-api/internal/domain"
)

// AtLeastTwo is a strategy that compares fields present in at least 2 products
type AtLeastTwo struct{}

// NewAtLeastTwo creates a new instance of the AtLeastTwo strategy
func NewAtLeastTwo() *AtLeastTwo {
	return &AtLeastTwo{}
}

// Name returns the identifier of this strategy
func (s *AtLeastTwo) Name() string {
	return "at_least_two"
}

// ResolveFields determines which fields to compare according to the rule "at least 2"
func (s *AtLeastTwo) ResolveFields(items []domain.Item, requested *[]string) []string {
	if len(items) == 0 {
		return []string{}
	}

	// STEP 1: Build the set of all possible fields of the items
	allFieldsMap := make(map[string]int) // field → count of items that have it

	for _, item := range items {
		// Root fields (always present in items of JSON)
		if item.Price != 0 || true { // Price always is
			allFieldsMap["price"]++
		}
		if item.Rating != 0 || true { // Rating always is
			allFieldsMap["rating"]++
		}

		// Specifications fields (nested)
		for specKey := range item.Specifications {
			fieldPath := fmt.Sprintf("specifications.%s", specKey)
			// Verify that the value is not nil before counting
			if value := item.Specifications[specKey]; value != nil {
				allFieldsMap[fieldPath]++
			}
		}
	}

	// STEP 2: Filter fields that appear in at least 2 items
	candidateFields := []string{}
	for field, count := range allFieldsMap {
		if count >= 2 {
			candidateFields = append(candidateFields, field)
		}
	}

	// STEP 3: If the client specified fields, do intersection
	var resolved []string
	if requested != nil && len(*requested) > 0 {
		// Create set of candidate fields for fast search
		candidateSet := make(map[string]bool)
		for _, f := range candidateFields {
			candidateSet[f] = true
		}

		// Iterate over the requested fields keeping the order
		for _, field := range *requested {
			if candidateSet[field] {
				resolved = append(resolved, field)
			}
		}
	} else {
		resolved = candidateFields
	}

	// Sort alphabetically
	sort.Strings(resolved)

	return resolved
}

// ComputeDiff calculates the differences for each resolved field
func (s *AtLeastTwo) ComputeDiff(ctx context.Context, items []domain.Item, resolved []string) (map[string]domain.DiffField, error) {
	diff := make(map[string]domain.DiffField)

	for _, fieldPath := range resolved {
		// Extract values of each item for this field
		values := make(map[string]interface{})
		for _, item := range items {
			val := s.extractFieldValue(item, fieldPath)
			values[item.ID] = val
		}

		// Determine metric for this field
		metric := GetMetricForField(fieldPath)

		// Calculate the best(s) according to the metric
		best := s.calculateBest(values, metric)

		diff[fieldPath] = domain.DiffField{
			Values: values,
			Metric: metric,
			Best:   best,
		}
	}

	return diff, nil
}

// extractFieldValue extracts the value of a field from an item.
// Supports root fields (e.g., "price") and nested fields (e.g., "specifications.buttons")
func (s *AtLeastTwo) extractFieldValue(item domain.Item, fieldPath string) interface{} {
	parts := strings.Split(fieldPath, ".")

	if len(parts) == 1 {
		// Campo root
		switch parts[0] {
		case "price":
			return item.Price
		case "rating":
			return item.Rating
		case "name":
			return item.Name
		case "description":
			return item.Description
		case "image_url":
			return item.ImageURL
		default:
			return nil
		}
	}

	if len(parts) == 2 && parts[0] == "specifications" {
		// Nested field in specifications
		specKey := parts[1]
		if val, exists := item.Specifications[specKey]; exists {
			// Extract the numeric value if it is an object with "value"
			if mapVal, ok := val.(map[string]interface{}); ok {
				if numVal, hasValue := mapVal["value"]; hasValue {
					return numVal
				}
			}
			return val
		}
		return nil
	}

	return nil
}

// calculateBest determines which items have the best value according to the metric
func (s *AtLeastTwo) calculateBest(values map[string]interface{}, metric *domain.Metric) []string {
	if metric == nil {
		return []string{}
	}

	// Filter values that are not nil
	validValues := make(map[string]interface{})
	for id, val := range values {
		if val != nil {
			validValues[id] = val
		}
	}

	if len(validValues) == 0 {
		return []string{}
	}

	best := []string{}

	switch *metric {
	case domain.LowerIsBetter:
		best = s.findLowest(validValues)
	case domain.HigherIsBetter:
		best = s.findHighest(validValues)
	case domain.TrueIsBetter:
		best = s.findTrueBest(validValues)
	}

	// Sort IDs for consistency
	sort.Strings(best)
	return best
}

// findLowest finds the IDs with the lowest numeric value
func (s *AtLeastTwo) findLowest(values map[string]interface{}) []string {
	var minVal *float64
	bestIDs := []string{}

	for id, val := range values {
		numVal := s.toFloat64(val)
		if numVal == nil {
			continue
		}

		if minVal == nil || *numVal < *minVal {
			minVal = numVal
			bestIDs = []string{id}
		} else if *numVal == *minVal {
			bestIDs = append(bestIDs, id)
		}
	}

	return bestIDs
}

// findHighest finds the IDs with the highest numeric value
func (s *AtLeastTwo) findHighest(values map[string]interface{}) []string {
	var maxVal *float64
	bestIDs := []string{}

	for id, val := range values {
		numVal := s.toFloat64(val)
		if numVal == nil {
			continue
		}

		if maxVal == nil || *numVal > *maxVal {
			maxVal = numVal
			bestIDs = []string{id}
		} else if *numVal == *maxVal {
			bestIDs = append(bestIDs, id)
		}
	}

	return bestIDs
}

// findTrueBest finds the IDs with the boolean value true
func (s *AtLeastTwo) findTrueBest(values map[string]interface{}) []string {
	bestIDs := []string{}

	for id, val := range values {
		boolVal := s.toBool(val)
		if boolVal != nil && *boolVal {
			bestIDs = append(bestIDs, id)
		}
	}

	return bestIDs
}

// toFloat64 converts an interface{} to a float64 if possible
func (s *AtLeastTwo) toFloat64(val interface{}) *float64 {
	if val == nil {
		return nil
	}

	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f := float64(v.Int())
		return &f
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f := float64(v.Uint())
		return &f
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		return &f
	default:
		return nil
	}
}

// toBool converts an interface{} to a bool if possible
func (s *AtLeastTwo) toBool(val interface{}) *bool {
	if val == nil {
		return nil
	}

	if b, ok := val.(bool); ok {
		return &b
	}

	return nil
}
