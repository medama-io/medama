package services

import (
	"net/http"

	"github.com/medama-io/medama/api"
)

// ErrBadRequest returns an API specific BadRequestError pointer.
func ErrBadRequest(err error) *api.BadRequestErrorHeaders {
	return &api.BadRequestErrorHeaders{
		Response: api.BadRequestError{
			Error: api.BadRequestErrorError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		},
	}
}

// ErrConflict returns an API specific ConflictError pointer.
func ErrConflict(err error) *api.ConflictErrorHeaders {
	return &api.ConflictErrorHeaders{
		Response: api.ConflictError{
			Error: api.ConflictErrorError{
				Code:    http.StatusConflict,
				Message: err.Error(),
			},
		},
	}
}

// ErrForbidden returns an API specific ForbiddenError pointer.
func ErrForbidden(err error) *api.ForbiddenErrorHeaders {
	return &api.ForbiddenErrorHeaders{
		Response: api.ForbiddenError{
			Error: api.ForbiddenErrorError{
				Code:    http.StatusForbidden,
				Message: err.Error(),
			},
		},
	}
}

// ErrNotFound returns an API specific NotFoundError pointer.
func ErrNotFound(err error) *api.NotFoundErrorHeaders {
	return &api.NotFoundErrorHeaders{
		Response: api.NotFoundError{
			Error: api.NotFoundErrorError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			},
		},
	}
}

// ErrInternalServerError returns an API specific InternalServerError pointer.
func ErrInternalServerError(err error) *api.InternalServerErrorHeaders {
	return &api.InternalServerErrorHeaders{
		Response: api.InternalServerError{
			Error: api.InternalServerErrorError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		},
	}
}

// ErrUnauthorised returns an API specific UnauthorisedError pointer.
func ErrUnauthorised(err error) *api.UnauthorisedErrorHeaders {
	return &api.UnauthorisedErrorHeaders{
		Response: api.UnauthorisedError{
			Error: api.UnauthorisedErrorError{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			},
		},
	}
}
