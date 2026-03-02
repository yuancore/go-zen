package zen

import "context"

// BaseModule provides a default (no-op) implementation of Module.
// Embed this in your module struct so you only need to override
// the methods you care about.
//
// Example:
//
//	type MyModule struct {
//	    zen.BaseModule
//	}
//
//	func NewMyModule() *MyModule {
//	    return &MyModule{
//	        BaseModule: zen.BaseModule{ModuleName: "my-module"},
//	    }
//	}
type BaseModule struct {
	ModuleName   string
	Dependencies []string
}

func (m *BaseModule) Name() string      { return m.ModuleName }
func (m *BaseModule) Depends() []string { return m.Dependencies }

func (m *BaseModule) Init(_ *App) error            { return nil }
func (m *BaseModule) Start() error                 { return nil }
func (m *BaseModule) Stop(_ context.Context) error { return nil }
