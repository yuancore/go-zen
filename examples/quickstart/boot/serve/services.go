package serve

import (
	"github.com/yuancore/go-zen/examples/quickstart/app/dao"
	"github.com/yuancore/go-zen/examples/quickstart/app/http/api"
	"github.com/yuancore/go-zen/examples/quickstart/app/service"
	"github.com/yuancore/go-zen/zen"
)

type services struct {
	systemAPI   *api.SystemAPI
	userAPI     *api.UserAPI
	userService *service.UserService
}

func newServices(app *zen.App) services {
	userDAO := dao.NewUserDAO(app)
	userService := service.NewUserService(userDAO, app.Logger())

	return services{
		systemAPI:   api.NewSystemAPI(app, app.Logger()),
		userAPI:     api.NewUserAPI(userService, app.Logger()),
		userService: userService,
	}
}
