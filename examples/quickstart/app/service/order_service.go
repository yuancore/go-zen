package service

import (
	"context"

	"github.com/yuancore/go-zen/examples/quickstart/app/dao"
	"github.com/yuancore/go-zen/examples/quickstart/app/entity/models"
	"github.com/yuancore/go-zen/examples/quickstart/app/entity/request"
)

// OrdersService handles business logic for the orders resource.
// It accepts a context.Context so it can be called from HTTP handlers,
// background jobs, CLI commands, or tests — all without holding *zen.App.
type OrdersService struct {
	dao *dao.OrdersDao
}

// NewOrdersService creates an OrdersService.
// ctx is the request/operation context — DB and Redis are resolved from it
// via the middleware injected at app startup (zdb.InjectMiddleware, zredis.InjectMiddleware).
func NewOrdersService(ctx context.Context) *OrdersService {
	return &OrdersService{
		dao: dao.NewOrdersDao(ctx),
	}
}

// Index returns a paginated list of orders.
func (svc *OrdersService) Index(req request.OrdersRequest) ([]models.Orders, int64, error) {
	return svc.dao.GetPage(req.PageParam)
}

// Show returns a single order by ID.
func (svc *OrdersService) Show(req request.OrdersRequest) (models.Orders, error) {
	return svc.dao.GetByID(req.Id)
}

// Store creates a new order from the form data.
func (svc *OrdersService) Store(req request.OrdersRequestForm) error {
	data := req.Orders
	_, err := svc.dao.Create(&data)
	return err
}

// Update saves changes to an existing order.
func (svc *OrdersService) Update(req request.OrdersRequestForm) error {
	return svc.dao.Update(req.Orders)
}

// Delete removes a single order by ID.
func (svc *OrdersService) Delete(req request.OrdersRequest) error {
	return svc.dao.DeleteByID(req.Id)
}

// Deletes removes multiple orders by their IDs.
func (svc *OrdersService) Deletes(req request.IdsRequest) error {
	return svc.dao.DeleteByIDs(req.Ids)
}
