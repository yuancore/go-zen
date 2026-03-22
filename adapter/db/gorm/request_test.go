package gormadapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yuancore/go-zen/zen"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestWithRequestBindsHTTPRequestContext(t *testing.T) {
	db := newMockGormDB(t)

	reqCtx := context.WithValue(context.Background(), requestContextKey("request_id"), "req-1")
	req := httptest.NewRequest(http.MethodGet, "/users", nil).WithContext(reqCtx)

	scoped := WithRequest(db, testZenContext{request: req})
	if scoped == db {
		t.Fatal("WithRequest should return a request-bound clone")
	}
	if scoped.Statement == nil {
		t.Fatal("WithRequest should populate statement context")
	}
	if got := scoped.Statement.Context.Value(requestContextKey("request_id")); got != "req-1" {
		t.Fatalf("statement context value = %v, want req-1", got)
	}
}

func TestScopeNamedBindsRequestContext(t *testing.T) {
	app := zen.New(zen.WithEngine(testEngine{}))
	app.Provide(DefaultServiceName, newMockGormDB(t))
	app.Provide(NamedService("analytics"), newMockGormDB(t))

	reqCtx := context.WithValue(context.Background(), requestContextKey("request_id"), "req-2")
	req := httptest.NewRequest(http.MethodGet, "/reports", nil).WithContext(reqCtx)

	scoped := For(app, testZenContext{request: req}).Named("analytics")
	if scoped.Statement == nil {
		t.Fatal("Named should return a request-bound db")
	}
	if got := scoped.Statement.Context.Value(requestContextKey("request_id")); got != "req-2" {
		t.Fatalf("statement context value = %v, want req-2", got)
	}
}

func TestConnectionBindsRequestContext(t *testing.T) {
	app := zen.New(zen.WithEngine(testEngine{}))
	app.Provide(DefaultServiceName, newMockGormDB(t))
	app.Provide(NamedService("mysql1"), newMockGormDB(t))

	reqCtx := context.WithValue(context.Background(), requestContextKey("request_id"), "req-3")
	req := httptest.NewRequest(http.MethodGet, "/users", nil).WithContext(reqCtx)

	db := Connection(app, testZenContext{request: req}, "mysql1")
	if db.Statement == nil {
		t.Fatal("Connection should return a request-bound db")
	}
	if got := db.Statement.Context.Value(requestContextKey("request_id")); got != "req-3" {
		t.Fatalf("statement context value = %v, want req-3", got)
	}
}

func TestTableOnUsesNamedConnectionAndBindsContext(t *testing.T) {
	app := zen.New(zen.WithEngine(testEngine{}))
	app.Provide(DefaultServiceName, newMockGormDB(t))
	app.Provide(NamedService("mysql1"), newMockGormDB(t))

	reqCtx := context.WithValue(context.Background(), requestContextKey("request_id"), "req-4")
	req := httptest.NewRequest(http.MethodGet, "/users", nil).WithContext(reqCtx)

	db := TableOn(app, testZenContext{request: req}, "mysql1", "sys_admin_users")
	if db.Statement == nil {
		t.Fatal("TableOn should return a table-scoped db")
	}
	if db.Statement.Table != "sys_admin_users" {
		t.Fatalf("statement table = %q, want sys_admin_users", db.Statement.Table)
	}
	if got := db.Statement.Context.Value(requestContextKey("request_id")); got != "req-4" {
		t.Fatalf("statement context value = %v, want req-4", got)
	}
}

func TestScopeResolveNamedMissingReturnsFalse(t *testing.T) {
	app := zen.New(zen.WithEngine(testEngine{}))

	db, ok := For(app, nil).ResolveNamed("missing")
	if ok {
		t.Fatal("ResolveNamed should report missing database")
	}
	if db != nil {
		t.Fatal("ResolveNamed should return nil for missing database")
	}
}

func newMockGormDB(t *testing.T) *gorm.DB {
	t.Helper()

	sqlDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	db, err := gorm.Open(gormmysql.New(gormmysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("open gorm db: %v", err)
	}
	return db
}

type requestContextKey string

type testZenContext struct {
	request *http.Request
}

func (c testZenContext) Param(string) string                 { return "" }
func (c testZenContext) Query(string) string                 { return "" }
func (c testZenContext) DefaultQuery(string, string) string  { return "" }
func (c testZenContext) Header(string) string                { return "" }
func (c testZenContext) SetHeader(string, string)            {}
func (c testZenContext) BindJSON(any) error                  { return nil }
func (c testZenContext) BindQuery(any) error                 { return nil }
func (c testZenContext) ShouldBind(any) error                { return nil }
func (c testZenContext) JSON(int, any)                       {}
func (c testZenContext) String(int, string, ...any)          {}
func (c testZenContext) Data(int, string, []byte)            {}
func (c testZenContext) Status(int)                          {}
func (c testZenContext) Redirect(int, string)                {}
func (c testZenContext) Request() *http.Request              { return c.request }
func (c testZenContext) ResponseWriter() http.ResponseWriter { return httptest.NewRecorder() }
func (c testZenContext) ClientIP() string                    { return "" }
func (c testZenContext) FullPath() string                    { return "" }
func (c testZenContext) Set(string, any)                     {}
func (c testZenContext) Get(string) (any, bool)              { return nil, false }
func (c testZenContext) MustGet(string) any                  { return nil }
func (c testZenContext) Next()                               {}
func (c testZenContext) Abort()                              {}
func (c testZenContext) AbortWithStatusJSON(int, any)        {}
func (c testZenContext) IsAborted() bool                     { return false }
