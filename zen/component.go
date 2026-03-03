package zen

import "context"

// Component is a pluggable unit of functionality with lifecycle hooks.
// This is the core building block for extending the framework.
//
// Lifecycle order:
//
//  1. Init()  — called in dependency order; use to read config, create connections
//  2. Start() — called in dependency order; use to start background tasks
//  3. Stop()  — called in reverse dependency order; use to close connections
//
// Example:
//
//	type RedisComponent struct {
//	    zen.BaseComponent
//	    client *redis.Client
//	}
//
//	func NewRedis() *RedisComponent {
//	    return &RedisComponent{
//	        BaseComponent: zen.BaseComponent{
//	            ComponentName: "redis",
//	            Deps:          []string{"config"},
//	        },
//	    }
//	}
//
//	func (r *RedisComponent) Init(app *App) error {
//	    addr := app.Config().GetString("redis.addr")
//	    r.client = redis.NewClient(&redis.Options{Addr: addr})
//	    app.Provide("redis", r.client)
//	    return nil
//	}
type Component interface {
	// Name returns the unique name of this component.
	Name() string

	// Init is called to initialize the component. Config and service
	// container are available via app.
	Init(app *App) error

	// Start is called after all components are initialized.
	Start() error

	// Stop is called during graceful shutdown (reverse order).
	Stop(ctx context.Context) error

	// Depends returns names of components this one depends on.
	Depends() []string
}

// BaseComponent provides a no-op implementation of Component.
// Embed it in your struct and override only the methods you need.
//
// Example:
//
//	type MyComp struct {
//	    zen.BaseComponent
//	}
//
//	func New() *MyComp {
//	    return &MyComp{
//	        BaseComponent: zen.BaseComponent{ComponentName: "mycomp"},
//	    }
//	}
type BaseComponent struct {
	ComponentName string
	Deps          []string
}

func (b *BaseComponent) Name() string                 { return b.ComponentName }
func (b *BaseComponent) Depends() []string            { return b.Deps }
func (b *BaseComponent) Init(_ *App) error            { return nil }
func (b *BaseComponent) Start() error                 { return nil }
func (b *BaseComponent) Stop(_ context.Context) error { return nil }

// Module is kept as an alias for backward compatibility.
// Deprecated: Use Component instead.
type Module = Component

// BaseModule is kept as an alias for backward compatibility.
// Deprecated: Use BaseComponent instead.
type BaseModule = BaseComponent
