package middleware

import (
	"net/http"

	"github.com/yuancore/go-zen/zen"
)

func SecurityHeaders() zen.Handler {
	return func(c zen.Context) {
		c.SetHeader("Cache-Control", "no-store")
		c.SetHeader("Pragma", "no-cache")
		c.SetHeader("X-Content-Type-Options", "nosniff")
		c.SetHeader("X-Frame-Options", "DENY")
		c.SetHeader("Referrer-Policy", "no-referrer")
		c.SetHeader("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.SetHeader("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'")
		c.Next()
	}
}

func MaxBodyBytes(limit int64) zen.Handler {
	return func(c zen.Context) {
		if limit > 0 {
			if request := c.Request(); request != nil {
				if request.ContentLength > limit {
					c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, map[string]any{
						"code":       http.StatusRequestEntityTooLarge,
						"message":    "request body too large",
						"request_id": requestID(c),
					})
					return
				}
				request.Body = http.MaxBytesReader(c.ResponseWriter(), request.Body, limit)
			}
		}
		c.Next()
	}
}

func requestID(c zen.Context) string {
	value, ok := c.Get("request_id")
	if !ok {
		return ""
	}
	requestID, _ := value.(string)
	return requestID
}
