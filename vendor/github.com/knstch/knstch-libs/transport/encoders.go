package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/knstch/knstch-libs/svcerrs"
)

func mapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, svcerrs.ErrDataNotFound):
		return http.StatusNotFound
	case errors.Is(err, svcerrs.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, svcerrs.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, svcerrs.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, svcerrs.ErrInvalidData):
		return http.StatusBadRequest
	case errors.Is(err, svcerrs.ErrGone):
		return http.StatusGone
	default:
		return http.StatusInternalServerError
	}
}

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	status := mapErrorToStatus(err)
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
