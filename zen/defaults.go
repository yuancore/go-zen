package zen

import (
	"log"
	"os"
)

// stdLogger is a fallback logger using standard library log.
type stdLogger struct {
	l      *log.Logger
	fields []any
}

func newStdLogger() *stdLogger {
	return &stdLogger{l: log.New(os.Stdout, "[zen] ", log.LstdFlags)}
}

func (s *stdLogger) Debug(msg string, kv ...any) { s.out("DBG", msg, kv) }
func (s *stdLogger) Info(msg string, kv ...any)  { s.out("INF", msg, kv) }
func (s *stdLogger) Warn(msg string, kv ...any)  { s.out("WRN", msg, kv) }
func (s *stdLogger) Error(msg string, kv ...any) { s.out("ERR", msg, kv) }
func (s *stdLogger) Fatal(msg string, kv ...any) { s.out("FTL", msg, kv); os.Exit(1) }

func (s *stdLogger) With(kv ...any) Logger {
	return &stdLogger{l: s.l, fields: append(append([]any{}, s.fields...), kv...)}
}

func (s *stdLogger) out(lvl, msg string, kv []any) {
	all := append(append([]any{}, s.fields...), kv...)
	if len(all) == 0 {
		s.l.Printf("%s %s", lvl, msg)
	} else {
		s.l.Printf("%s %s %v", lvl, msg, all)
	}
}

// emptyConfig is a no-op Config used when no config is provided.
type emptyConfig struct{}

func (emptyConfig) GetString(string) string            { return "" }
func (emptyConfig) GetInt(string) int                  { return 0 }
func (emptyConfig) GetBool(string) bool                { return false }
func (emptyConfig) GetFloat64(string) float64          { return 0 }
func (emptyConfig) GetStringSlice(string) []string     { return nil }
func (emptyConfig) GetStringMap(string) map[string]any { return nil }
func (e emptyConfig) Sub(string) Config                { return e }
func (emptyConfig) Unmarshal(string, any) error        { return nil }
func (emptyConfig) IsSet(string) bool                  { return false }
func (emptyConfig) Set(string, any)                    {}
