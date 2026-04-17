package router

import (
	"fmt"

	"github.com/yuancore/go-zen/examples/quickstart/app/http/api"
	"github.com/yuancore/go-zen/examples/quickstart/routes"
	"github.com/yuancore/go-zen/zen"
)

// Setup wires up all routes. It is called from serve.OnStart so that
// the database component is already initialized in the container.
func Setup(app *zen.App) error {
	if app == nil {
		return fmt.Errorf("register routes: nil app")
	}

	// Global middleware
	app.Middleware()

	// Controllers — no *App needed; DB/Cache are resolved per-request via context
	sysCtrl := api.NewSystemController()
	ordersCtrl := api.NewOrdersController()

	// Health / liveness
	app.GET(routes.Ping, sysCtrl.Ping)

	// Orders CRUD
	v1 := app.Group(routes.APIV1)
	{
		v1.GET(routes.Orders, ordersCtrl.Index)
		v1.GET(routes.OrderByID, ordersCtrl.Show)
		v1.POST(routes.Orders, ordersCtrl.Create)
		v1.PUT(routes.OrderByID, ordersCtrl.Update)
		v1.DELETE(routes.Orders, ordersCtrl.Delete)
	}

	return nil
}
