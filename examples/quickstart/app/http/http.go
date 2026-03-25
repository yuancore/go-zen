package httpapp

import "github.com/yuancore/go-zen/zen"

const DefaultRequestBodyLimitBytes int64 = 1 << 20

func RequestBodyLimitBytes(cfg zen.Config) int64 {
	if cfg == nil {
		return DefaultRequestBodyLimitBytes
	}
	limit := cfg.GetInt("system.request_body_limit_bytes")
	if limit <= 0 {
		return DefaultRequestBodyLimitBytes
	}
	return int64(limit)
}
