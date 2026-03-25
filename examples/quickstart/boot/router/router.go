package router

import (
	"fmt"

	httpapp "github.com/yuancore/go-zen/examples/quickstart/app/http"
	"github.com/yuancore/go-zen/examples/quickstart/app/http/api"
	"github.com/yuancore/go-zen/examples/quickstart/app/http/middleware"
	"github.com/yuancore/go-zen/examples/quickstart/routes"
	"github.com/yuancore/go-zen/zen"
)

func Register(app *zen.App, systemAPI *api.SystemAPI, userAPI *api.UserAPI) error {
	switch {
	case app == nil:
		return fmt.Errorf("register routes: nil app")
	case systemAPI == nil:
		return fmt.Errorf("register routes: nil system api")
	case userAPI == nil:
		return fmt.Errorf("register routes: nil user api")
	}

	app.Middleware(
		middleware.SecurityHeaders(),
		middleware.MaxBodyBytes(httpapp.RequestBodyLimitBytes(app.Config())),
	)

	app.GET(routes.Ping, systemAPI.Ping)

	apiV1 := app.Group(routes.APIV1)
	apiV1.GET(routes.Users, userAPI.List)
	apiV1.GET(routes.UserByID, userAPI.Get)
	apiV1.POST(routes.Users, userAPI.Create)
	apiV1.PUT(routes.UserByID, userAPI.Update)
	apiV1.DELETE(routes.UserByID, userAPI.Delete)

	return nil
}
