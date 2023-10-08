package services

import (
	"net/http"

	"github.com/medama-io/medama/api"
)

// ErrBadRequest returns an API specific BadRequestError pointer.
func ErrBadRequest(err error) *api.BadRequestError {
	return &api.BadRequestError{
		Error: api.BadRequestErrorError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		},
	}
}

// ErrConflict returns an API specific ConflictError pointer.
func ErrConflict(err error) *api.ConflictError {
	return &api.ConflictError{
		Error: api.ConflictErrorError{
			Code:    http.StatusConflict,
			Message: err.Error(),
		},
	}
}

// ErrNotFound returns an API specific NotFoundError pointer.
func ErrNotFound(err error) *api.NotFoundError {
	return &api.NotFoundError{
		Error: api.NotFoundErrorError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		},
	}
}

// ErrInternalServerError returns an API specific InternalServerError pointer.
func ErrInternalServerError(err error) *api.InternalServerError {
	return &api.InternalServerError{
		Error: api.InternalServerErrorError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

// ErrUnauthorised returns an API specific UnauthorisedError pointer.
func ErrUnauthorised(err error) *api.UnauthorisedError {
	return &api.UnauthorisedError{
		Error: api.UnauthorisedErrorError{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		},
	}
}
