package serve

import (
	"context"
	"fmt"

	"github.com/yuancore/go-zen/examples/quickstart/app/service"
	"github.com/yuancore/go-zen/zen"
)

func registerStartupHooks(app *zen.App, userService *service.UserService) {
	app.OnStart(func() error {
		if err := userService.Migrate(context.Background()); err != nil {
			return fmt.Errorf("migrate user schema: %w", err)
		}
		return nil
	})
}
