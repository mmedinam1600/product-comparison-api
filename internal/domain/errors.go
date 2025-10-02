package domain

import "net/http"

// ErrorCode representa los códigos de error de la API
type ErrorCode string

const (
	ErrorCodeAtLeastTwoIds  ErrorCode = "AtLeastTwoIds"
	ErrorCodeIdNotFound     ErrorCode = "IdNotFound"
	ErrorCodeUnknownField   ErrorCode = "UnknownField"
	ErrorCodeMissingField   ErrorCode = "MissingField"
	ErrorCodeInvalidRequest ErrorCode = "InvalidRequest"
	ErrorCodeConflict       ErrorCode = "Conflict"
)

// ErrorResponse representa la respuesta de error de la API
type ErrorResponse struct {
	ErrorCode     ErrorCode `json:"error_code"`
	Message       string    `json:"message"`
	MissingIDs    []string  `json:"missing_ids,omitempty"`
	UnknownFields []string  `json:"unknown_fields,omitempty"`
}

// HTTPStatusCode retorna el código HTTP apropiado para cada error
func (e ErrorCode) HTTPStatusCode() int {
	switch e {
	case ErrorCodeIdNotFound:
		return http.StatusNotFound
	case ErrorCodeAtLeastTwoIds, ErrorCodeUnknownField:
		return http.StatusUnprocessableEntity
	case ErrorCodeMissingField, ErrorCodeInvalidRequest:
		return http.StatusBadRequest
	case ErrorCodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
