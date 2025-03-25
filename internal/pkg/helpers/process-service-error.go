package helpers

import (
	"calendar-api/internal/core"
	internalerrors "calendar-api/internal/pkg/errors"
	"errors"
	"github.com/go-chi/render"
	"net/http"
)

func ProcessServiceError(w http.ResponseWriter, r *http.Request, err error) (statusCode int, message string) {
	var APIError internalerrors.APIError
	switch {
	case errors.As(err, &APIError):
		statusCode = APIError.StatusCode()
		message = APIError.Error()
	default:
		statusCode = http.StatusInternalServerError
		message = err.Error()
	}

	w.WriteHeader(statusCode)
	render.JSON(w, r, core.Response[any]{
		Ok:    false,
		Error: message,
	})

	return statusCode, message
}
