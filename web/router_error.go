package web

import (
	"errors"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	message := err.Error()

	var e *Error
	switch errors.Unwrap(err) {
	case ErrBadRequest:
		e = NewError(message, http.StatusBadRequest)
	case ErrWeakPassword:
		e = NewError(message, http.StatusBadRequest)
	case ErrUnprocessableEntity:
		e = NewError(message, http.StatusUnprocessableEntity)
	case ErrResourceNotFound:
		e = NewError(message, http.StatusNotFound)
	case ErrInvalidToken:
		e = NewError(message, http.StatusForbidden)
	case ErrAlteredTokenClaims:
		e = NewError(message, http.StatusForbidden)
	case ErrResourceAlreadyExists:
		e = NewError(message, http.StatusConflict)
	default:
		e = NewError(message, http.StatusInternalServerError)
	}

	_ = RespondJSON(w, e, e.StatusCode)
}
