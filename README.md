# go-zen

## GORM Component

`adapter/db/gorm` provides a production-oriented database component for `zen.App`:

- No package-level global DB state
- Multi-connection registry with default DB injection
- Connection pool tuning and startup ping
- GORM SQL logging wired into the framework logger
- Graceful close during app shutdown

Example:

```go
package main

import (
	gormadapter "github.com/yuancore/go-zen/adapter/db/gorm"
)

dbModule := gormadapter.New()

app := zen.New(
	zen.WithConfig(cfg),
	zen.WithLogger(logger),
	zen.WithEngine(engine),
)

app.Use(dbModule)

db := gormadapter.MustResolve(app)
analytics := gormadapter.MustResolveNamed(app, "analytics")

_ = db
_ = analytics
```

Recommended config:

```toml
[database]
default = "main"
prepare_stmt = true
skip_default_transaction = true
ping_timeout = 3

[database.logger]
enabled = true
level = "warn"
slow_threshold_millis = 200
ignore_record_not_found_error = true

[[database.connections]]
name = "main"
driver = "mysql"
host = "127.0.0.1"
port = 3306
username = "app"
password = "secret"
database = "app"
params = "charset=utf8mb4&parseTime=True&loc=Local"
max_idle_conns = 20
max_open_conns = 100
conn_max_lifetime = 1800
conn_max_idle_time = 300
log = true
level = "warn"
```

Legacy root-level `[[connections]]` config is also supported for compatibility with antgo-style layout.

