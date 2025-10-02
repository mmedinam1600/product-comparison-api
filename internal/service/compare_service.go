package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/mmedinam1600/product-comparison-api/internal/data"
	"github.com/mmedinam1600/product-comparison-api/internal/domain"
	"github.com/mmedinam1600/product-comparison-api/internal/service/strategy"
	"go.uber.org/zap"
)

// CompareService define the contract for the comparison service
type CompareService interface {
	// Compare executes the product comparison
	Compare(ctx context.Context, req domain.CompareRequest) (domain.CompareResult, domain.Metadata, *domain.ErrorResponse)

	// GenerateCacheKey generates a unique key for caching based on the IDs
	GenerateCacheKey(ids []string) string
}

// CompareServiceImpl implements CompareService
type CompareServiceImpl struct {
	repo       data.CatalogRepository
	strategies map[string]strategy.Interface
	logger     *zap.Logger
}

// NewCompareService creates a new instance of the service
func NewCompareService(repo data.CatalogRepository, logger *zap.Logger) *CompareServiceImpl {
	// Register available strategies
	strategies := make(map[string]strategy.Interface)
	atLeastTwo := strategy.NewAtLeastTwo()
	strategies[atLeastTwo.Name()] = atLeastTwo

	return &CompareServiceImpl{
		repo:       repo,
		strategies: strategies,
		logger:     logger,
	}
}

// Compare implements CompareService.Compare
func (s *CompareServiceImpl) Compare(ctx context.Context, req domain.CompareRequest) (domain.CompareResult, domain.Metadata, *domain.ErrorResponse) {
	// === STEP 1: Validate IDs ===
	// Validate that there are at least 2 unique IDs
	uniqueIDs := s.getUniqueIDs(req.Ids)
	if len(uniqueIDs) < 2 {
		return domain.CompareResult{}, domain.Metadata{}, &domain.ErrorResponse{
			ErrorCode: domain.ErrorCodeAtLeastTwoIds,
			Message:   "At least 2 unique ids are required.",
		}
	}

	s.logger.Debug("validating IDs", zap.Int("unique_count", len(uniqueIDs)))

	// === STEP 2: Resolve items from the repository ===
	items, missingIDs := s.repo.GetByIDs(ctx, uniqueIDs)

	if len(missingIDs) > 0 {
		s.logger.Warn("some IDs not found", zap.Strings("missing", missingIDs))
		return domain.CompareResult{}, domain.Metadata{}, &domain.ErrorResponse{
			ErrorCode:  domain.ErrorCodeIdNotFound,
			Message:    "Some products were not found.",
			MissingIDs: missingIDs,
		}
	}

	s.logger.Debug("items resolved", zap.Int("count", len(items)))

	// === STEP 3: Select and apply strategy ===
	// For now, we only use "at_least_two"
	strategyName := "at_least_two"
	strat, exists := s.strategies[strategyName]
	if !exists {
		s.logger.Error("strategy not found", zap.String("strategy", strategyName))
		return domain.CompareResult{}, domain.Metadata{}, &domain.ErrorResponse{
			ErrorCode: domain.ErrorCodeInvalidRequest,
			Message:   fmt.Sprintf("Strategy '%s' not available.", strategyName),
		}
	}

	// === STEP 4: Resolve comparable fields ===
	resolvedFields := strat.ResolveFields(items, req.Fields)

	s.logger.Debug("fields resolved",
		zap.Int("resolved_count", len(resolvedFields)),
		zap.Strings("resolved", resolvedFields),
	)

	// If fields were specified but none are valid → error
	if req.Fields != nil && len(*req.Fields) > 0 && len(resolvedFields) == 0 {
		return domain.CompareResult{}, domain.Metadata{}, &domain.ErrorResponse{
			ErrorCode:     domain.ErrorCodeUnknownField,
			Message:       "Unknown fields requested.",
			UnknownFields: *req.Fields,
		}
	}

	// If there are no comparable fields at all
	if len(resolvedFields) == 0 {
		return domain.CompareResult{}, domain.Metadata{}, &domain.ErrorResponse{
			ErrorCode: domain.ErrorCodeUnknownField,
			Message:   "No comparable fields found.",
		}
	}

	// === STEP 5: Calculate differences ===
	diff, err := strat.ComputeDiff(ctx, items, resolvedFields)
	if err != nil {
		s.logger.Error("failed to compute diff", zap.Error(err))
		return domain.CompareResult{}, domain.Metadata{}, &domain.ErrorResponse{
			ErrorCode: domain.ErrorCodeInvalidRequest,
			Message:   "Failed to compute differences.",
		}
	}

	// === STEP 6: Calculate comparability score ===
	// baseCandidate: if the client sent fields → use those; if not → use resolvedFields
	var baseCandidate []string
	if req.Fields != nil && len(*req.Fields) > 0 {
		baseCandidate = *req.Fields
	} else {
		baseCandidate = resolvedFields
	}

	comparabilityScore := 0.0
	if len(baseCandidate) > 0 {
		comparabilityScore = float64(len(resolvedFields)) / float64(len(baseCandidate))
	}

	s.logger.Debug("comparability calculated",
		zap.Float64("score", comparabilityScore),
		zap.Int("base_count", len(baseCandidate)),
	)

	// === STEP 7: Build result and metadata ===
	result := domain.CompareResult{
		Items:        items,
		SharedFields: resolvedFields,
		Diff:         diff,
	}

	metadata := domain.Metadata{
		Order:           uniqueIDs,
		RequestedFields: req.Fields,
		ResolvedFields:  resolvedFields,
		ComparePolicy: domain.ComparePolicy{
			EffectiveMode:      strat.Name(),
			ComparabilityScore: comparabilityScore,
		},
		Currency: "USD",
		Version:  "1.0",
	}

	s.logger.Info("comparison completed successfully",
		zap.Int("items", len(items)),
		zap.Int("fields", len(resolvedFields)),
	)

	return result, metadata, nil
}

// GenerateCacheKey generates a cache key based on the ordered IDs
func (s *CompareServiceImpl) GenerateCacheKey(ids []string) string {
	// Get unique and ordered IDs
	uniqueIDs := s.getUniqueIDs(ids)
	sort.Strings(uniqueIDs)

	// Generate SHA-256 hash
	joinedIDs := strings.Join(uniqueIDs, ",")
	hash := sha256.Sum256([]byte(joinedIDs))
	return hex.EncodeToString(hash[:])
}

// getUniqueIDs returns a list of unique IDs maintaining the original order
func (s *CompareServiceImpl) getUniqueIDs(ids []string) []string {
	seen := make(map[string]bool)
	unique := make([]string, 0, len(ids))

	for _, id := range ids {
		if id != "" && !seen[id] {
			seen[id] = true
			unique = append(unique, id)
		}
	}

	return unique
}
