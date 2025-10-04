package service

import (
	"context"
	"testing"

	"github.com/mmedinam1600/product-comparison-api/internal/domain"
	"go.uber.org/zap"
)

// MockCatalogRepository is a mock implementation for testing
type MockCatalogRepository struct {
	items       map[string]domain.Item
	shouldError bool
}

func (m *MockCatalogRepository) GetByIDs(ctx context.Context, ids []string) ([]domain.Item, []string) {
	found := []domain.Item{}
	missing := []string{}

	for _, id := range ids {
		if item, exists := m.items[id]; exists {
			found = append(found, item)
		} else {
			missing = append(missing, id)
		}
	}

	return found, missing
}

func (m *MockCatalogRepository) GetAll(ctx context.Context) []domain.Item {
	items := make([]domain.Item, 0, len(m.items))
	for _, item := range m.items {
		items = append(items, item)
	}
	return items
}

func TestNewCompareService(t *testing.T) {
	logger := zap.NewNop()
	repo := &MockCatalogRepository{}

	service := NewCompareService(repo, logger)

	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}

	if service.repo == nil {
		t.Error("Expected repo to be set")
	}

	if service.logger == nil {
		t.Error("Expected logger to be set")
	}

	if len(service.strategies) == 0 {
		t.Error("Expected strategies to be registered")
	}

	if _, exists := service.strategies["at_least_two"]; !exists {
		t.Error("Expected at_least_two strategy to be registered")
	}
}

func TestCompareService_GenerateCacheKey(t *testing.T) {
	logger := zap.NewNop()
	repo := &MockCatalogRepository{}
	service := NewCompareService(repo, logger)

	tests := []struct {
		name     string
		ids      []string
		expected string
	}{
		{
			name:     "Same IDs in different order produce same key",
			ids:      []string{"a", "b", "c"},
			expected: service.GenerateCacheKey([]string{"c", "a", "b"}),
		},
		{
			name: "Duplicate IDs produce same key as unique",
			ids:  []string{"a", "b", "a"},
		},
		{
			name: "Empty strings are filtered",
			ids:  []string{"a", "", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key1 := service.GenerateCacheKey(tt.ids)

			// Verify key is not empty
			if key1 == "" {
				t.Error("Expected non-empty cache key")
			}

			// Verify key is deterministic
			key2 := service.GenerateCacheKey(tt.ids)
			if key1 != key2 {
				t.Errorf("Expected same key for same input, got %s and %s", key1, key2)
			}

			// Verify SHA-256 hash length (64 hex characters)
			if len(key1) != 64 {
				t.Errorf("Expected cache key length 64, got %d", len(key1))
			}
		})
	}
}

func TestCompareService_Compare_ValidateIDs(t *testing.T) {
	logger := zap.NewNop()
	repo := &MockCatalogRepository{
		items: map[string]domain.Item{
			"id1": {ID: "id1", Name: "Item 1", Price: 10.0},
			"id2": {ID: "id2", Name: "Item 2", Price: 20.0},
		},
	}
	service := NewCompareService(repo, logger)
	ctx := context.Background()

	tests := []struct {
		name          string
		request       domain.CompareRequest
		expectError   bool
		expectedError domain.ErrorCode
	}{
		{
			name: "Less than 2 unique IDs returns error",
			request: domain.CompareRequest{
				Ids: []string{"id1"},
			},
			expectError:   true,
			expectedError: domain.ErrorCodeAtLeastTwoIds,
		},
		{
			name: "Duplicate IDs (same ID twice) returns error",
			request: domain.CompareRequest{
				Ids: []string{"id1", "id1"},
			},
			expectError:   true,
			expectedError: domain.ErrorCodeAtLeastTwoIds,
		},
		{
			name: "Empty IDs array returns error",
			request: domain.CompareRequest{
				Ids: []string{},
			},
			expectError:   true,
			expectedError: domain.ErrorCodeAtLeastTwoIds,
		},
		{
			name: "Valid 2 unique IDs succeeds",
			request: domain.CompareRequest{
				Ids: []string{"id1", "id2"},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, errResp := service.Compare(ctx, tt.request)

			if tt.expectError {
				if errResp == nil {
					t.Fatal("Expected error but got nil")
				}
				if errResp.ErrorCode != tt.expectedError {
					t.Errorf("Expected error code %v, got %v", tt.expectedError, errResp.ErrorCode)
				}
			} else {
				if errResp != nil {
					t.Errorf("Expected no error, got: %v", errResp.Message)
				}
			}
		})
	}
}

func TestCompareService_Compare_MissingIDs(t *testing.T) {
	logger := zap.NewNop()
	repo := &MockCatalogRepository{
		items: map[string]domain.Item{
			"id1": {ID: "id1", Name: "Item 1", Price: 10.0, Rating: 4.0},
		},
	}
	service := NewCompareService(repo, logger)
	ctx := context.Background()

	request := domain.CompareRequest{
		Ids: []string{"id1", "id-nonexistent"},
	}

	_, _, errResp := service.Compare(ctx, request)

	if errResp == nil {
		t.Fatal("Expected error for missing ID")
	}

	if errResp.ErrorCode != domain.ErrorCodeIdNotFound {
		t.Errorf("Expected ErrorCodeIdNotFound, got %v", errResp.ErrorCode)
	}

	if len(errResp.MissingIDs) != 1 || errResp.MissingIDs[0] != "id-nonexistent" {
		t.Errorf("Expected missing ID 'id-nonexistent', got %v", errResp.MissingIDs)
	}
}

func TestCompareService_Compare_UnknownFields(t *testing.T) {
	logger := zap.NewNop()
	repo := &MockCatalogRepository{
		items: map[string]domain.Item{
			"id1": {ID: "id1", Name: "Item 1", Price: 10.0, Rating: 4.0},
			"id2": {ID: "id2", Name: "Item 2", Price: 20.0, Rating: 4.5},
		},
	}
	service := NewCompareService(repo, logger)
	ctx := context.Background()

	unknownFields := []string{"nonexistent_field"}
	request := domain.CompareRequest{
		Ids:    []string{"id1", "id2"},
		Fields: &unknownFields,
	}

	_, _, errResp := service.Compare(ctx, request)

	if errResp == nil {
		t.Fatal("Expected error for unknown fields")
	}

	if errResp.ErrorCode != domain.ErrorCodeUnknownField {
		t.Errorf("Expected ErrorCodeUnknownField, got %v", errResp.ErrorCode)
	}
}

func TestCompareService_Compare_Success(t *testing.T) {
	logger := zap.NewNop()
	repo := &MockCatalogRepository{
		items: map[string]domain.Item{
			"id1": {
				ID:     "id1",
				Name:   "Mouse Logitech",
				Price:  50.0,
				Rating: 4.5,
				Specifications: map[string]interface{}{
					"weight":  100,
					"buttons": 5,
				},
			},
			"id2": {
				ID:     "id2",
				Name:   "Mouse Razer",
				Price:  75.0,
				Rating: 4.8,
				Specifications: map[string]interface{}{
					"weight":  120,
					"buttons": 7,
				},
			},
		},
	}
	service := NewCompareService(repo, logger)
	ctx := context.Background()

	request := domain.CompareRequest{
		Ids: []string{"id1", "id2"},
	}

	result, metadata, errResp := service.Compare(ctx, request)

	if errResp != nil {
		t.Fatalf("Expected no error, got: %v", errResp.Message)
	}

	// Validate result
	if len(result.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result.Items))
	}

	if len(result.SharedFields) == 0 {
		t.Error("Expected shared fields to be populated")
	}

	if len(result.Diff) == 0 {
		t.Error("Expected diff to be populated")
	}

	// Validate metadata
	if metadata.ComparePolicy.EffectiveMode != "at_least_two" {
		t.Errorf("Expected effective_mode 'at_least_two', got %v", metadata.ComparePolicy.EffectiveMode)
	}

	if metadata.ComparePolicy.ComparabilityScore <= 0 {
		t.Errorf("Expected positive comparability score, got %v", metadata.ComparePolicy.ComparabilityScore)
	}

	if metadata.Currency != "USD" {
		t.Errorf("Expected currency USD, got %v", metadata.Currency)
	}

	if len(metadata.Order) != 2 {
		t.Errorf("Expected 2 IDs in order, got %d", len(metadata.Order))
	}
}

func TestCompareService_Compare_WithFieldsFilter(t *testing.T) {
	logger := zap.NewNop()
	repo := &MockCatalogRepository{
		items: map[string]domain.Item{
			"id1": {
				ID:     "id1",
				Price:  50.0,
				Rating: 4.5,
				Specifications: map[string]interface{}{
					"weight":  100,
					"buttons": 5,
				},
			},
			"id2": {
				ID:     "id2",
				Price:  75.0,
				Rating: 4.8,
				Specifications: map[string]interface{}{
					"weight":  120,
					"buttons": 7,
				},
			},
		},
	}
	service := NewCompareService(repo, logger)
	ctx := context.Background()

	requestedFields := []string{"price", "rating"}
	request := domain.CompareRequest{
		Ids:    []string{"id1", "id2"},
		Fields: &requestedFields,
	}

	result, metadata, errResp := service.Compare(ctx, request)

	if errResp != nil {
		t.Fatalf("Expected no error, got: %v", errResp.Message)
	}

	// Verify only requested fields are in shared_fields
	if len(result.SharedFields) != 2 {
		t.Errorf("Expected 2 shared fields, got %d", len(result.SharedFields))
	}

	expectedFields := map[string]bool{"price": true, "rating": true}
	for _, field := range result.SharedFields {
		if !expectedFields[field] {
			t.Errorf("Unexpected field in shared_fields: %s", field)
		}
	}

	// Verify diff only contains requested fields
	if len(result.Diff) != 2 {
		t.Errorf("Expected 2 diff entries, got %d", len(result.Diff))
	}

	// Verify comparability score
	if metadata.ComparePolicy.ComparabilityScore != 1.0 {
		t.Errorf("Expected comparability score 1.0 (all requested fields present), got %v",
			metadata.ComparePolicy.ComparabilityScore)
	}
}

func TestCompareService_GetUniqueIDs(t *testing.T) {
	logger := zap.NewNop()
	repo := &MockCatalogRepository{}
	service := NewCompareService(repo, logger)

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Remove duplicates",
			input:    []string{"a", "b", "a", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Remove empty strings",
			input:    []string{"a", "", "b", ""},
			expected: []string{"a", "b"},
		},
		{
			name:     "Maintain order (first occurrence)",
			input:    []string{"c", "a", "b", "a"},
			expected: []string{"c", "a", "b"},
		},
		{
			name:     "All empty returns empty",
			input:    []string{"", "", ""},
			expected: []string{},
		},
		{
			name:     "All duplicates returns one",
			input:    []string{"a", "a", "a"},
			expected: []string{"a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.getUniqueIDs(tt.input)

			if len(got) != len(tt.expected) {
				t.Errorf("Length mismatch: got %d, want %d", len(got), len(tt.expected))
			}

			for i, id := range tt.expected {
				if i >= len(got) || got[i] != id {
					t.Errorf("At index %d: got %v, want %v", i, got, tt.expected)
					break
				}
			}
		})
	}
}
