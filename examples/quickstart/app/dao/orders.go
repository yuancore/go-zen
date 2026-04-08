package dao

import (
	"github.com/yuancore/go-zen/examples/quickstart/app/entity/models"
	"github.com/yuancore/zentool/page"
	"gorm.io/gorm"
)

// OrdersDao encapsulates all database operations for the orders table.
// The *gorm.DB passed in should already have the request context attached
// (e.g. via gormadapter.DB(app, ctx)).
type OrdersDao struct {
	db *gorm.DB
}

func NewOrdersDao(db *gorm.DB) *OrdersDao {
	return &OrdersDao{db: db}
}

// Create inserts a new order and returns its generated primary key.
func (d *OrdersDao) Create(order *models.Orders) (int, error) {
	if err := d.db.Create(order).Error; err != nil {
		return 0, err
	}
	return order.Id, nil
}

// DeleteByID soft-deletes the order with the given ID.
func (d *OrdersDao) DeleteByID(id int) error {
	return d.db.Delete(&models.Orders{}, id).Error
}

// DeleteByIDs soft-deletes orders whose IDs are in the given slice.
func (d *OrdersDao) DeleteByIDs(ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.Delete(&models.Orders{}, ids).Error
}

// Update saves non-zero fields of the given order (must have primary key set).
func (d *OrdersDao) Update(order models.Orders) error {
	return d.db.Updates(&order).Error
}

// GetList returns all orders (no pagination).
func (d *OrdersDao) GetList() ([]models.Orders, error) {
	var list []models.Orders
	err := d.db.Model(&models.Orders{}).Find(&list).Error
	return list, err
}

// GetPage returns a paginated list and the total row count.
func (d *OrdersDao) GetPage(p page.PageParam) (list []models.Orders, total int64, err error) {
	tx := d.db.Model(&models.Orders{})

	if err = tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (p.CurrentPage - 1) * p.PageSize
	if offset < 0 {
		offset = 0
	}
	size := p.PageSize
	if size <= 0 {
		size = 10
	}

	err = tx.Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

// GetByID returns the order with the given ID, or a zero-value struct if not found.
func (d *OrdersDao) GetByID(id int) (models.Orders, error) {
	var row models.Orders
	err := d.db.Where("id = ?", id).First(&row).Error
	return row, err
}
