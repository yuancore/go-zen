package request

import (
	"github.com/yuancore/go-zen/examples/quickstart/app/entity/models"
	"github.com/yuancore/zentool/page"
)

type OrdersRequest struct {
	models.Orders
	page.PageParam
}

type OrdersRequestForm struct {
	models.Orders
}
