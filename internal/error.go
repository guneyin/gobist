package internal

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidHttpMethod = errors.New("INVALID_HTTP_METHOD")
	ErrInvalidPayload    = errors.New("INVALID_PAYLOAD")
)

func ErrHttpError(r *http.Response) error {
	return fmt.Errorf("HTTP_ERROR: STATUS=%d, MESSAGE=%s", r.StatusCode, r.Status)
}
