package strategy

import (
	"context"
	"reflect"
	"testing"

	"github.com/mmedinam1600/product-comparison-api/internal/domain"
)

func TestAtLeastTwo_Name(t *testing.T) {
	strategy := NewAtLeastTwo()
	expected := "at_least_two"

	if strategy.Name() != expected {
		t.Errorf("Name() = %v, want %v", strategy.Name(), expected)
	}
}

func TestAtLeastTwo_ResolveFields(t *testing.T) {
	strategy := NewAtLeastTwo()

	tests := []struct {
		name      string
		items     []domain.Item
		requested *[]string
		expected  []string
	}{
		{
			name:      "Empty items returns empty",
			items:     []domain.Item{},
			requested: nil,
			expected:  []string{},
		},
		{
			name: "Two items with common fields",
			items: []domain.Item{
				{
					ID:     "1",
					Price:  10.0,
					Rating: 4.5,
					Specifications: map[string]interface{}{
						"weight":  100,
						"buttons": 5,
					},
				},
				{
					ID:     "2",
					Price:  20.0,
					Rating: 4.0,
					Specifications: map[string]interface{}{
						"weight":  120,
						"buttons": 3,
					},
				},
			},
			requested: nil,
			expected:  []string{"price", "rating", "specifications.buttons", "specifications.weight"},
		},
		{
			name: "Field present in only one item is excluded",
			items: []domain.Item{
				{
					ID:     "1",
					Price:  10.0,
					Rating: 4.5,
					Specifications: map[string]interface{}{
						"weight":  100,
						"buttons": 5,
					},
				},
				{
					ID:     "2",
					Price:  20.0,
					Rating: 4.0,
					Specifications: map[string]interface{}{
						"weight": 120,
						// No "buttons"
					},
				},
			},
			requested: nil,
			expected:  []string{"price", "rating", "specifications.weight"},
		},
		{
			name: "Requested fields filters correctly",
			items: []domain.Item{
				{
					ID:     "1",
					Price:  10.0,
					Rating: 4.5,
					Specifications: map[string]interface{}{
						"weight":  100,
						"buttons": 5,
					},
				},
				{
					ID:     "2",
					Price:  20.0,
					Rating: 4.0,
					Specifications: map[string]interface{}{
						"weight":  120,
						"buttons": 3,
					},
				},
			},
			requested: &[]string{"price", "rating"},
			expected:  []string{"price", "rating"},
		},
		{
			name: "Requested field not in candidate returns empty",
			items: []domain.Item{
				{
					ID:     "1",
					Price:  10.0,
					Rating: 4.5,
				},
				{
					ID:     "2",
					Price:  20.0,
					Rating: 4.0,
				},
			},
			requested: &[]string{"specifications.nonexistent"},
			expected:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strategy.ResolveFields(tt.items, tt.requested)

			// Handle nil vs empty slice
			if len(got) == 0 && len(tt.expected) == 0 {
				return // Both are empty, pass
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ResolveFields() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAtLeastTwo_ComputeDiff(t *testing.T) {
	strategy := NewAtLeastTwo()
	ctx := context.Background()

	tests := []struct {
		name     string
		items    []domain.Item
		resolved []string
		validate func(t *testing.T, diff map[string]domain.DiffField)
	}{
		{
			name: "Compute diff for price (lower_is_better)",
			items: []domain.Item{
				{ID: "1", Price: 100.0},
				{ID: "2", Price: 50.0},
				{ID: "3", Price: 75.0},
			},
			resolved: []string{"price"},
			validate: func(t *testing.T, diff map[string]domain.DiffField) {
				priceDiff, exists := diff["price"]
				if !exists {
					t.Fatal("price diff not found")
				}

				if priceDiff.Metric == nil {
					t.Fatal("metric is nil")
				}

				if *priceDiff.Metric != domain.LowerIsBetter {
					t.Errorf("Expected LowerIsBetter metric, got %v", *priceDiff.Metric)
				}

				if len(priceDiff.Best) != 1 || priceDiff.Best[0] != "2" {
					t.Errorf("Expected best = [2], got %v", priceDiff.Best)
				}

				if priceDiff.Values["2"] != 50.0 {
					t.Errorf("Expected value 50.0 for ID 2, got %v", priceDiff.Values["2"])
				}
			},
		},
		{
			name: "Compute diff for rating (higher_is_better)",
			items: []domain.Item{
				{ID: "1", Rating: 4.0},
				{ID: "2", Rating: 4.8},
				{ID: "3", Rating: 4.2},
			},
			resolved: []string{"rating"},
			validate: func(t *testing.T, diff map[string]domain.DiffField) {
				ratingDiff, exists := diff["rating"]
				if !exists {
					t.Fatal("rating diff not found")
				}

				if *ratingDiff.Metric != domain.HigherIsBetter {
					t.Errorf("Expected HigherIsBetter metric, got %v", *ratingDiff.Metric)
				}

				if len(ratingDiff.Best) != 1 || ratingDiff.Best[0] != "2" {
					t.Errorf("Expected best = [2], got %v", ratingDiff.Best)
				}
			},
		},
		{
			name: "Multiple items with same best value",
			items: []domain.Item{
				{ID: "1", Price: 50.0},
				{ID: "2", Price: 50.0},
				{ID: "3", Price: 75.0},
			},
			resolved: []string{"price"},
			validate: func(t *testing.T, diff map[string]domain.DiffField) {
				priceDiff := diff["price"]

				// Both "1" and "2" should be best
				if len(priceDiff.Best) != 2 {
					t.Errorf("Expected 2 best items, got %d", len(priceDiff.Best))
				}

				expectedBest := map[string]bool{"1": true, "2": true}
				for _, id := range priceDiff.Best {
					if !expectedBest[id] {
						t.Errorf("Unexpected best ID: %s", id)
					}
				}
			},
		},
		{
			name: "Boolean field (true_is_better)",
			items: []domain.Item{
				{
					ID: "1",
					Specifications: map[string]interface{}{
						"wireless": true,
					},
				},
				{
					ID: "2",
					Specifications: map[string]interface{}{
						"wireless": false,
					},
				},
			},
			resolved: []string{"specifications.wireless"},
			validate: func(t *testing.T, diff map[string]domain.DiffField) {
				wirelessDiff := diff["specifications.wireless"]

				if *wirelessDiff.Metric != domain.TrueIsBetter {
					t.Errorf("Expected TrueIsBetter metric, got %v", *wirelessDiff.Metric)
				}

				if len(wirelessDiff.Best) != 1 || wirelessDiff.Best[0] != "1" {
					t.Errorf("Expected best = [1], got %v", wirelessDiff.Best)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, err := strategy.ComputeDiff(ctx, tt.items, tt.resolved)
			if err != nil {
				t.Fatalf("ComputeDiff() error = %v", err)
			}

			tt.validate(t, diff)
		})
	}
}

func TestAtLeastTwo_ExtractFieldValue(t *testing.T) {
	strategy := NewAtLeastTwo()

	item := domain.Item{
		ID:          "test-1",
		Name:        "Test Product",
		Description: "Description",
		Price:       99.99,
		Rating:      4.5,
		ImageURL:    "http://example.com/image.jpg",
		Specifications: map[string]interface{}{
			"weight":   500,
			"buttons":  7,
			"wireless": true,
		},
	}

	tests := []struct {
		name      string
		fieldPath string
		expected  interface{}
	}{
		{
			name:      "Extract root field: price",
			fieldPath: "price",
			expected:  99.99,
		},
		{
			name:      "Extract root field: rating",
			fieldPath: "rating",
			expected:  4.5,
		},
		{
			name:      "Extract root field: name",
			fieldPath: "name",
			expected:  "Test Product",
		},
		{
			name:      "Extract nested field: specifications.weight",
			fieldPath: "specifications.weight",
			expected:  500,
		},
		{
			name:      "Extract nested field: specifications.buttons",
			fieldPath: "specifications.buttons",
			expected:  7,
		},
		{
			name:      "Extract nested boolean: specifications.wireless",
			fieldPath: "specifications.wireless",
			expected:  true,
		},
		{
			name:      "Nonexistent root field returns nil",
			fieldPath: "nonexistent",
			expected:  nil,
		},
		{
			name:      "Nonexistent nested field returns nil",
			fieldPath: "specifications.nonexistent",
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strategy.extractFieldValue(item, tt.fieldPath)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("extractFieldValue() = %v (%T), want %v (%T)",
					got, got, tt.expected, tt.expected)
			}
		})
	}
}
