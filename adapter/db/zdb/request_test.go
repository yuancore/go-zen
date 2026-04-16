package zdb

import (
	"context"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

// testZenContext is a minimal zen.Context stub used only in tests.
// Only Request() is meaningful; all other methods return zero values.
type testZenContext struct {
	request *http.Request
}

// --- context.Context ---
func (c testZenContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (c testZenContext) Done() <-chan struct{}       { return nil }
func (c testZenContext) Err() error                  { return nil }
func (c testZenContext) Value(any) any               { return nil }

// --- URL params ---
func (c testZenContext) Param(string) string                          { return "" }
func (c testZenContext) AddParam(string, string)                      {}
func (c testZenContext) Query(string) string                          { return "" }
func (c testZenContext) DefaultQuery(string, string) string           { return "" }
func (c testZenContext) GetQuery(string) (string, bool)               { return "", false }
func (c testZenContext) QueryArray(string) []string                   { return nil }
func (c testZenContext) GetQueryArray(string) ([]string, bool)        { return nil, false }
func (c testZenContext) QueryMap(string) map[string]string            { return nil }
func (c testZenContext) GetQueryMap(string) (map[string]string, bool) { return nil, false }

// --- Request info ---
func (c testZenContext) Header(string) string                { return "" }
func (c testZenContext) SetHeader(string, string)            {}
func (c testZenContext) ContentType() string                 { return "" }
func (c testZenContext) GetRawData() ([]byte, error)         { return nil, nil }
func (c testZenContext) IsWebsocket() bool                   { return false }
func (c testZenContext) RemoteIP() string                    { return "" }
func (c testZenContext) ClientIP() string                    { return "" }
func (c testZenContext) FullPath() string                    { return "" }
func (c testZenContext) Request() *http.Request              { return c.request }
func (c testZenContext) ResponseWriter() http.ResponseWriter { return httptest.NewRecorder() }
func (c testZenContext) HandlerName() string                 { return "" }
func (c testZenContext) HandlerNames() []string              { return nil }

// --- Bind ---
func (c testZenContext) Bind(any) error       { return nil }
func (c testZenContext) BindJSON(any) error   { return nil }
func (c testZenContext) BindXML(any) error    { return nil }
func (c testZenContext) BindQuery(any) error  { return nil }
func (c testZenContext) BindYAML(any) error   { return nil }
func (c testZenContext) BindTOML(any) error   { return nil }
func (c testZenContext) BindPlain(any) error  { return nil }
func (c testZenContext) BindUri(any) error    { return nil }
func (c testZenContext) BindHeader(any) error { return nil }

// --- ShouldBind ---
func (c testZenContext) ShouldBind(any) error              { return nil }
func (c testZenContext) ShouldBindJSON(any) error          { return nil }
func (c testZenContext) ShouldBindXML(any) error           { return nil }
func (c testZenContext) ShouldBindQuery(any) error         { return nil }
func (c testZenContext) ShouldBindYAML(any) error          { return nil }
func (c testZenContext) ShouldBindTOML(any) error          { return nil }
func (c testZenContext) ShouldBindPlain(any) error         { return nil }
func (c testZenContext) ShouldBindUri(any) error           { return nil }
func (c testZenContext) ShouldBindHeader(any) error        { return nil }
func (c testZenContext) ShouldBindBodyWithJSON(any) error  { return nil }
func (c testZenContext) ShouldBindBodyWithXML(any) error   { return nil }
func (c testZenContext) ShouldBindBodyWithYAML(any) error  { return nil }
func (c testZenContext) ShouldBindBodyWithTOML(any) error  { return nil }
func (c testZenContext) ShouldBindBodyWithPlain(any) error { return nil }

// --- Form ---
func (c testZenContext) PostForm(string) string                          { return "" }
func (c testZenContext) DefaultPostForm(string, string) string           { return "" }
func (c testZenContext) GetPostForm(string) (string, bool)               { return "", false }
func (c testZenContext) PostFormArray(string) []string                   { return nil }
func (c testZenContext) GetPostFormArray(string) ([]string, bool)        { return nil, false }
func (c testZenContext) PostFormMap(string) map[string]string            { return nil }
func (c testZenContext) GetPostFormMap(string) (map[string]string, bool) { return nil, false }

// --- Multipart ---
func (c testZenContext) FormFile(string) (*multipart.FileHeader, error) { return nil, nil }
func (c testZenContext) MultipartForm() (*multipart.Form, error)        { return nil, nil }
func (c testZenContext) SaveUploadedFile(*multipart.FileHeader, string, ...fs.FileMode) error {
	return nil
}

// --- Cookie ---
func (c testZenContext) Cookie(string) (string, error)                             { return "", nil }
func (c testZenContext) SetCookie(string, string, int, string, string, bool, bool) {}
func (c testZenContext) SetCookieData(*http.Cookie)                                {}
func (c testZenContext) SetSameSite(http.SameSite)                                 {}

// --- Response ---
func (c testZenContext) JSON(int, any)                                                   {}
func (c testZenContext) SecureJSON(int, any)                                             {}
func (c testZenContext) JSONP(int, any)                                                  {}
func (c testZenContext) IndentedJSON(int, any)                                           {}
func (c testZenContext) AsciiJSON(int, any)                                              {}
func (c testZenContext) PureJSON(int, any)                                               {}
func (c testZenContext) HTML(int, string, any)                                           {}
func (c testZenContext) XML(int, any)                                                    {}
func (c testZenContext) YAML(int, any)                                                   {}
func (c testZenContext) TOML(int, any)                                                   {}
func (c testZenContext) ProtoBuf(int, any)                                               {}
func (c testZenContext) BSON(int, any)                                                   {}
func (c testZenContext) String(int, string, ...any)                                      {}
func (c testZenContext) Data(int, string, []byte)                                        {}
func (c testZenContext) DataFromReader(int, int64, string, io.Reader, map[string]string) {}
func (c testZenContext) Status(int)                                                      {}
func (c testZenContext) Redirect(int, string)                                            {}
func (c testZenContext) File(string)                                                     {}
func (c testZenContext) FileFromFS(string, http.FileSystem)                              {}
func (c testZenContext) FileAttachment(string, string)                                   {}
func (c testZenContext) SSEvent(string, any)                                             {}
func (c testZenContext) Stream(func(io.Writer) bool) bool                                { return false }
func (c testZenContext) NegotiateFormat(...string) string                                { return "" }
func (c testZenContext) SetAccepted(...string)                                           {}

// --- Context store ---
func (c testZenContext) Set(string, any)                                    {}
func (c testZenContext) Get(string) (any, bool)                             { return nil, false }
func (c testZenContext) MustGet(string) any                                 { return nil }
func (c testZenContext) Delete(string)                                      {}
func (c testZenContext) GetString(string) string                            { return "" }
func (c testZenContext) GetBool(string) bool                                { return false }
func (c testZenContext) GetInt(string) int                                  { return 0 }
func (c testZenContext) GetInt8(string) int8                                { return 0 }
func (c testZenContext) GetInt16(string) int16                              { return 0 }
func (c testZenContext) GetInt32(string) int32                              { return 0 }
func (c testZenContext) GetInt64(string) int64                              { return 0 }
func (c testZenContext) GetUint(string) uint                                { return 0 }
func (c testZenContext) GetUint8(string) uint8                              { return 0 }
func (c testZenContext) GetUint16(string) uint16                            { return 0 }
func (c testZenContext) GetUint32(string) uint32                            { return 0 }
func (c testZenContext) GetUint64(string) uint64                            { return 0 }
func (c testZenContext) GetFloat32(string) float32                          { return 0 }
func (c testZenContext) GetFloat64(string) float64                          { return 0 }
func (c testZenContext) GetTime(string) time.Time                           { return time.Time{} }
func (c testZenContext) GetDuration(string) time.Duration                   { return 0 }
func (c testZenContext) GetStringSlice(string) []string                     { return nil }
func (c testZenContext) GetStringMap(string) map[string]any                 { return nil }
func (c testZenContext) GetStringMapString(string) map[string]string        { return nil }
func (c testZenContext) GetStringMapStringSlice(string) map[string][]string { return nil }

// --- Error tracking & flow control ---
func (c testZenContext) AddError(error)                   {}
func (c testZenContext) Next()                            {}
func (c testZenContext) Abort()                           {}
func (c testZenContext) AbortWithStatus(int)              {}
func (c testZenContext) AbortWithStatusJSON(int, any)     {}
func (c testZenContext) AbortWithStatusPureJSON(int, any) {}
func (c testZenContext) AbortWithError(int, error)        {}
func (c testZenContext) IsAborted() bool                  { return false }
