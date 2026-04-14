package zen

import "time"

// Option configures App during construction via New().
type Option func(*App)

// WithName sets the application name (used in logs and health checks).
func WithName(name string) Option {
	return func(a *App) { a.name = name }
}

// WithConfig sets a custom Config implementation.
func WithConfig(c Config) Option {
	return func(a *App) { a.cfg = c }
}

// WithLogger sets a custom Logger implementation.
func WithLogger(l Logger) Option {
	return func(a *App) { a.logger = l }
}

// WithEngine sets the HTTP engine implementation.
func WithEngine(e Engine) Option {
	return func(a *App) { a.engine = e }
}

// StopTimeout sets the graceful-shutdown deadline.
func WithStopTimeout(d time.Duration) Option {
	return func(a *App) { a.stopTimeout = d }
}

// Banner sets a custom startup banner. Pass empty string to disable.
func Banner(b string) Option {
	return func(a *App) { a.banner = &b }
}
