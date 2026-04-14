package service

import (
	"github.com/yuancore/go-zen/examples/quickstart/app/dao"
	"github.com/yuancore/go-zen/examples/quickstart/app/entity/models"
	"github.com/yuancore/go-zen/examples/quickstart/app/entity/request"
	"gorm.io/gorm"
)

// OrdersService handles business logic for the orders resource.
// The *gorm.DB passed in should already carry the request context.
type OrdersService struct {
	dao *dao.OrdersDao
}

// NewOrdersService creates an OrdersService bound to the given database handle.
func NewOrdersService(db *gorm.DB) *OrdersService {
	return &OrdersService{
		dao: dao.NewOrdersDao(db),
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
