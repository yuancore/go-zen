package zdb

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yuancore/go-zen/zen"
	gormlogger "gorm.io/gorm/logger"
)

func TestTableName(t *testing.T) {
	cases := map[string]string{
		"SELECT * FROM `sys_admin_users` LIMIT 20":         "sys_admin_users",
		"insert into orders(id, name) values (?, ?)":       "orders",
		"UPDATE [user_profiles] SET name = ? WHERE id = ?": "user_profiles",
		"delete from \"audit_logs\" where id = $1":         "audit_logs",
		"select 1": "",
	}

	for sql, want := range cases {
		if got := tableName(sql); got != want {
			t.Fatalf("tableName(%q) = %q, want %q", sql, got, want)
		}
	}
}

func TestTraceUsesStructuredFields(t *testing.T) {
	rec := &recordLogger{}
	l := newSQLLogger(rec, resolvedLogConfig{
		enabled:                   true,
		level:                     gormlogger.Info,
		slowThreshold:             200 * time.Millisecond,
		ignoreRecordNotFoundError: true,
	}).(*sqlLogger)

	ctx := context.WithValue(context.Background(), "request_id", "req-1")
	begin := time.Now().Add(-1500 * time.Microsecond)
	l.Trace(ctx, begin, func() (string, int64) {
		return "SELECT * FROM `sys_admin_users` LIMIT 20", 18
	}, nil)

	if rec.msg != "db.query" {
		t.Fatalf("msg = %q, want db.query", rec.msg)
	}
	if rec.fields["event"] != "db.query" {
		t.Fatalf("event = %v, want db.query", rec.fields["event"])
	}
	if rec.fields["table"] != "sys_admin_users" {
		t.Fatalf("table = %v, want sys_admin_users", rec.fields["table"])
	}
	if rec.fields["request_id"] != "req-1" {
		t.Fatalf("request_id = %v, want req-1", rec.fields["request_id"])
	}
	if _, ok := rec.fields["duration_ms"]; !ok {
		t.Fatal("duration_ms missing")
	}
}

func TestTraceErrorUsesDbErrorEvent(t *testing.T) {
	rec := &recordLogger{}
	l := newSQLLogger(rec, resolvedLogConfig{
		enabled:                   true,
		level:                     gormlogger.Info,
		slowThreshold:             200 * time.Millisecond,
		ignoreRecordNotFoundError: true,
	}).(*sqlLogger)

	l.Trace(context.Background(), time.Now(), func() (string, int64) {
		return "SELECT * FROM `sys_admin_users`", 0
	}, errors.New("boom"))

	if rec.msg != "db.error" {
		t.Fatalf("msg = %q, want db.error", rec.msg)
	}
	if rec.fields["event"] != "db.error" {
		t.Fatalf("event = %v, want db.error", rec.fields["event"])
	}
}

type recordLogger struct {
	msg    string
	fields map[string]any
}

func (r *recordLogger) Debug(msg string, kv ...any) {}
func (r *recordLogger) Warn(msg string, kv ...any)  { r.capture(msg, kv...) }
func (r *recordLogger) Error(msg string, kv ...any) { r.capture(msg, kv...) }
func (r *recordLogger) Fatal(msg string, kv ...any) { r.capture(msg, kv...) }
func (r *recordLogger) Info(msg string, kv ...any)  { r.capture(msg, kv...) }
func (r *recordLogger) With(kv ...any) zen.Logger   { return r }

func (r *recordLogger) capture(msg string, kv ...any) {
	r.msg = msg
	r.fields = make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			continue
		}
		r.fields[key] = kv[i+1]
	}
}
