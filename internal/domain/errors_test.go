package domain

import (
	"net/http"
	"testing"
)

func TestErrorCode_HTTPStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected int
	}{
		{
			name:     "AtLeastTwoIds returns 422",
			code:     ErrorCodeAtLeastTwoIds,
			expected: http.StatusUnprocessableEntity,
		},
		{
			name:     "IdNotFound returns 404",
			code:     ErrorCodeIdNotFound,
			expected: http.StatusNotFound,
		},
		{
			name:     "UnknownField returns 422",
			code:     ErrorCodeUnknownField,
			expected: http.StatusUnprocessableEntity,
		},
		{
			name:     "MissingField returns 400",
			code:     ErrorCodeMissingField,
			expected: http.StatusBadRequest,
		},
		{
			name:     "InvalidRequest returns 400",
			code:     ErrorCodeInvalidRequest,
			expected: http.StatusBadRequest,
		},
		{
			name:     "Conflict returns 409",
			code:     ErrorCodeConflict,
			expected: http.StatusConflict,
		},
		{
			name:     "Unknown error code returns 500",
			code:     ErrorCode("UNKNOWN_CODE"),
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.code.HTTPStatusCode()
			if got != tt.expected {
				t.Errorf("HTTPStatusCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}
