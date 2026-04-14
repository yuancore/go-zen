package zdb

import (
	"context"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yuancore/go-zen/zen"
)

func TestComponentInitRegistersDefaultAndNamedConnections(t *testing.T) {
	primarySQLDB, primaryMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("create primary sqlmock: %v", err)
	}
	defer func() { _ = primarySQLDB.Close() }()

	analyticsSQLDB, analyticsMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("create analytics sqlmock: %v", err)
	}
	defer func() { _ = analyticsSQLDB.Close() }()

	primaryMock.ExpectPing()
	primaryMock.ExpectExec("select 1").WillReturnResult(sqlmock.NewResult(0, 1))
	primaryMock.ExpectClose()
	analyticsMock.ExpectPing()
	analyticsMock.ExpectExec("select 1").WillReturnResult(sqlmock.NewResult(0, 1))
	analyticsMock.ExpectClose()

	settings := DefaultSettings()
	settings.Default = "primary"
	settings.PrepareStmt = boolPtr(false)
	settings.Connections = []ConnectionConfig{
		{
			Name:   "primary",
			Driver: "mysql",
			sqlDB:  primarySQLDB,
		},
		{
			Name:       "analytics",
			Driver:     "mysql",
			sqlDB:      analyticsSQLDB,
			LogEnabled: boolPtr(false),
		},
	}

	component := New(WithSettings(settings))
	app := zen.New(zen.WithEngine(testEngine{}))

	if err := component.Init(app); err != nil {
		t.Fatalf("init component: %v", err)
	}

	db := MustResolve(app)
	if err := db.Exec("select 1").Error; err != nil {
		t.Fatalf("exec on default db: %v", err)
	}

	analytics := MustResolveNamed(app, "analytics")
	if err := analytics.Exec("select 1").Error; err != nil {
		t.Fatalf("exec on named db: %v", err)
	}

	manager := MustResolveManager(app)
	if got := manager.DefaultName(); got != "primary" {
		t.Fatalf("default name = %q, want primary", got)
	}
	if got := len(manager.Names()); got != 2 {
		t.Fatalf("connection count = %d, want 2", got)
	}

	if err := component.Stop(context.Background()); err != nil {
		t.Fatalf("stop component: %v", err)
	}
	if err := primaryMock.ExpectationsWereMet(); err != nil {
		t.Fatalf("primary expectations: %v", err)
	}
	if err := analyticsMock.ExpectationsWereMet(); err != nil {
		t.Fatalf("analytics expectations: %v", err)
	}
}

func TestComponentInitWithoutConnectionsSkips(t *testing.T) {
	component := New(WithSettings(DefaultSettings()))
	app := zen.New(zen.WithEngine(testEngine{}))

	if err := component.Init(app); err != nil {
		t.Fatalf("init component: %v", err)
	}

	if _, ok := Resolve(app); ok {
		t.Fatal("default db should not be registered")
	}
	if _, ok := ResolveManager(app); ok {
		t.Fatal("manager should not be registered")
	}
}

type testEngine struct{}

func (testEngine) GET(string, ...zen.Handler)                   {}
func (testEngine) POST(string, ...zen.Handler)                  {}
func (testEngine) PUT(string, ...zen.Handler)                   {}
func (testEngine) DELETE(string, ...zen.Handler)                {}
func (testEngine) PATCH(string, ...zen.Handler)                 {}
func (testEngine) HEAD(string, ...zen.Handler)                  {}
func (testEngine) OPTIONS(string, ...zen.Handler)               {}
func (testEngine) Any(string, ...zen.Handler)                   {}
func (testEngine) Handle(string, string, ...zen.Handler)        {}
func (testEngine) Use(...zen.Handler)                           {}
func (testEngine) Group(string, ...zen.Handler) zen.RouterGroup { return testRouterGroup{} }
func (testEngine) StaticFile(string, string)                    {}
func (testEngine) Static(string, string)                        {}
func (testEngine) StaticFS(string, http.FileSystem)             {}
func (testEngine) Start(string) error                           { return nil }
func (testEngine) Stop(context.Context) error                   { return nil }

type testRouterGroup struct{}

func (testRouterGroup) GET(string, ...zen.Handler)                   {}
func (testRouterGroup) POST(string, ...zen.Handler)                  {}
func (testRouterGroup) PUT(string, ...zen.Handler)                   {}
func (testRouterGroup) DELETE(string, ...zen.Handler)                {}
func (testRouterGroup) PATCH(string, ...zen.Handler)                 {}
func (testRouterGroup) HEAD(string, ...zen.Handler)                  {}
func (testRouterGroup) OPTIONS(string, ...zen.Handler)               {}
func (testRouterGroup) Any(string, ...zen.Handler)                   {}
func (testRouterGroup) Handle(string, string, ...zen.Handler)        {}
func (testRouterGroup) Use(...zen.Handler)                           {}
func (testRouterGroup) Group(string, ...zen.Handler) zen.RouterGroup { return testRouterGroup{} }
func (testRouterGroup) StaticFile(string, string)                    {}
func (testRouterGroup) Static(string, string)                        {}
func (testRouterGroup) StaticFS(string, http.FileSystem)             {}

var _ zen.Engine = testEngine{}
var _ zen.RouterGroup = testRouterGroup{}
