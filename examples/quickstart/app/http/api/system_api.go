package api

import (
	"net/http"

	zdb "github.com/yuancore/go-zen/adapter/db/gorm"
	"github.com/yuancore/go-zen/zen"
)

type SystemAPI struct {
	app    *zen.App
	logger zen.Logger
}

func NewSystemAPI(app *zen.App, logger zen.Logger) *SystemAPI {
	return &SystemAPI{
		app:    app,
		logger: logger.With("module", "system_api"),
	}
}

func (a *SystemAPI) Ping(c zen.Context) {
	if db, ok := zdb.Resolve(a.app); ok {
		sqlDB, err := db.DB()
		if err != nil {
			a.logger.Error("resolve sql db failed", "err", err)
			writeError(c, http.StatusServiceUnavailable, "database unavailable")
			return
		}
		if err := sqlDB.PingContext(c.Request().Context()); err != nil {
			a.logger.Error("ping database failed", "err", err)
			writeError(c, http.StatusServiceUnavailable, "database unavailable")
			return
		}
	}

	writeOK(c, map[string]any{
		"app":  a.app.AppName(),
		"pong": "ok",
	})
}
