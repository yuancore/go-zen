package api

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/yuancore/go-zen/zen"
)

type responseEnvelope struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
	Data      any    `json:"data,omitempty"`
}

func writeOK(c zen.Context, data any) {
	writeJSON(c, http.StatusOK, "ok", data)
}

func writeCreated(c zen.Context, data any) {
	writeJSON(c, http.StatusCreated, "created", data)
}

func writeError(c zen.Context, status int, message string) {
	writeJSON(c, status, message, nil)
}

func writeJSON(c zen.Context, status int, message string, data any) {
	c.JSON(status, responseEnvelope{
		Code:      status,
		Message:   message,
		RequestID: requestID(c),
		Data:      data,
	})
}

func bindJSON(c zen.Context, target any) bool {
	if err := c.BindJSON(target); err != nil {
		writeError(c, bindStatus(err), bindMessage(err))
		return false
	}
	return true
}

func bindQuery(c zen.Context, target any) bool {
	if err := c.BindQuery(target); err != nil {
		writeError(c, bindStatus(err), "invalid query parameters")
		return false
	}
	return true
}

func parseUint64ID(c zen.Context, key string) (uint64, bool) {
	value := strings.TrimSpace(c.Param(key))
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		writeError(c, http.StatusBadRequest, "invalid resource id")
		return 0, false
	}
	return id, true
}

func bindStatus(err error) int {
	switch {
	case err == nil:
		return http.StatusBadRequest
	case strings.Contains(strings.ToLower(err.Error()), "request body too large"):
		return http.StatusRequestEntityTooLarge
	default:
		return http.StatusBadRequest
	}
}

func bindMessage(err error) string {
	switch {
	case err == nil:
		return "invalid request"
	case errors.Is(err, io.EOF):
		return "request body is required"
	case strings.Contains(strings.ToLower(err.Error()), "request body too large"):
		return "request body too large"
	default:
		return "invalid request payload"
	}
}

func requestID(c zen.Context) string {
	if c == nil {
		return ""
	}
	value, ok := c.Get("request_id")
	if !ok {
		return ""
	}
	requestID, _ := value.(string)
	return requestID
}
