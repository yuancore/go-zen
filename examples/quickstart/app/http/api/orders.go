package api

import (
	"github.com/yuancore/go-zen/examples/quickstart/app/entity/request"
	"github.com/yuancore/go-zen/examples/quickstart/app/entity/vo"
	"github.com/yuancore/go-zen/examples/quickstart/app/service"
	"github.com/yuancore/go-zen/zen"
	"github.com/yuancore/zentool/page"
)

// OrdersController handles HTTP requests for the orders resource.
// It holds no *App reference — all dependencies are propagated via zen.Context.
type OrdersController struct {
	vo.Base
}

// NewOrdersController creates a controller. No app injection required.
func NewOrdersController() *OrdersController {
	return &OrdersController{}
}

// @Tags 订单
// @Summary 分页查询订单
// @Accept json
// @Produce json
// @Param data query page.PageParam true "分页参数"
// @Success 200 {object} response.Write{data=response.Page{items=[]models.Orders}}
// @Router /api/v1/orders [get]
func (ctrl *OrdersController) Index(c zen.Context) {
	req := request.OrdersRequest{PageParam: page.New()}
	if err := c.ShouldBindQuery(&req); err != nil {
		ctrl.Fail(c, vo.INVALID_REQUEST_PARAMETERS, err)
		return
	}

	list, total, err := service.NewOrdersService(c).Index(req)
	if err != nil {
		ctrl.Fail(c, vo.FAILED, err)
		return
	}
	ctrl.Success(c, vo.SUCCESS, ctrl.Page(total, list))
}

// @Tags 订单
// @Summary 查询单个订单
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} response.Write{data=models.Orders}
// @Router /api/v1/orders/:id [get]
func (ctrl *OrdersController) Show(c zen.Context) {
	var req request.OrdersRequest
	if err := c.ShouldBindUri(&req); err != nil {
		ctrl.Fail(c, vo.INVALID_REQUEST_PARAMETERS, err)
		return
	}

	result, err := service.NewOrdersService(c).Show(req)
	if err != nil {
		ctrl.Fail(c, vo.FAILED, err)
		return
	}
	ctrl.Success(c, vo.SUCCESS, result)
}

// @Tags 订单
// @Summary 创建订单
// @Accept json
// @Produce json
// @Param data body request.OrdersRequestForm true "创建参数"
// @Success 200 {object} response.Write
// @Router /api/v1/orders [post]
func (ctrl *OrdersController) Create(c zen.Context) {
	var req request.OrdersRequestForm
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Fail(c, vo.INVALID_REQUEST_PARAMETERS, err)
		return
	}

	if err := service.NewOrdersService(c).Store(req); err != nil {
		ctrl.Fail(c, vo.CREATION_FAILED, err)
		return
	}
	ctrl.Success(c, vo.CREATION_SUCCESS)
}

// @Tags 订单
// @Summary 更新订单
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param data body request.OrdersRequestForm true "更新参数"
// @Success 200 {object} response.Write
// @Router /api/v1/orders/:id [put]
func (ctrl *OrdersController) Update(c zen.Context) {
	var req request.OrdersRequestForm
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Fail(c, vo.INVALID_REQUEST_PARAMETERS, err)
		return
	}

	if err := service.NewOrdersService(c).Update(req); err != nil {
		ctrl.Fail(c, vo.UPDATE_FAILED, err)
		return
	}
	ctrl.Success(c, vo.UPDATE_SUCCESS)
}

// @Tags 订单
// @Summary 批量删除订单
// @Accept json
// @Produce json
// @Param data body request.IdsRequest true "删除参数"
// @Success 200 {object} response.Write
// @Router /api/v1/orders [delete]
func (ctrl *OrdersController) Delete(c zen.Context) {
	var req request.IdsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Fail(c, vo.INVALID_REQUEST_PARAMETERS, err)
		return
	}

	if err := service.NewOrdersService(c).Deletes(req); err != nil {
		ctrl.Fail(c, vo.DELETE_FAILED, err)
		return
	}
	ctrl.Success(c, vo.DELETE_SUCCESS)
}
