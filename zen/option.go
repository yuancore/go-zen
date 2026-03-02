package zen

import "time"

// Option configures App during construction via New().
type Option func(*App)

// WithEngine sets the HTTP engine adapter.
func WithEngine(e Engine) Option {
	return func(a *App) { a.engine = e }
}

// WithConfig sets the configuration backend.
func WithConfig(c Config) Option {
	return func(a *App) { a.cfg = c }
}

// WithLogger sets the structured logger.
func WithLogger(l Logger) Option {
	return func(a *App) { a.logger = l }
}

// WithStopTimeout sets the graceful-shutdown deadline.
func WithStopTimeout(d time.Duration) Option {
	return func(a *App) { a.stopTimeout = d }
}
