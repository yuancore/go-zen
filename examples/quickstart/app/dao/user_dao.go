package dao

import (
	"context"

	"github.com/yuancore/go-zen/examples/quickstart/app/entity/models"
	"gorm.io/gorm"
)

type CardsDao struct {
	db     *gorm.DB
	models *models.User
}

func NewCardsDao(ctx context.Context, db *gorm.DB) *CardsDao {
	if db == nil {
		db = ant.Db()
	}
	return &CardsDao{db: db.WithContext(ctx)}
}

// Create
func (dao *CardsDao) Create(cards *models.User) (int, error) {
	if err := dao.db.Create(cards).Error; err != nil {
		return 0, err
	}
	return cards.Id, nil
}

// DeleteById
func (dao *CardsDao) DeleteById(id int) error {
	return dao.db.Delete(&dao.models, id).Error
}

// DeleteByIds
func (dao *CardsDao) DeleteByIds(id []int) error {
	return dao.db.Delete(&dao.models, id).Error
}

// Update
func (dao *CardsDao) Update(cards models.User) error {
	return dao.db.Updates(&cards).Error
}

// GetList
func (dao *CardsDao) GetList() (list []models.User) {
	dao.db.Model(&dao.models).Find(&list)
	return list
}

// GetPage
func (dao *CardsDao) GetPage(page page.PageParam) (list []models.Cards, total int64, err error) {
	err = dao.db.Model(&dao.models).Scopes(
		asql.Where("project_id", "=", page.FilterMap["project_id"]), asql.Where("name", "LIKE", "%"+conv.String(page.FilterMap["name"])+"%"), asql.Where("is_pay_later", "=", page.FilterMap["is_pay_later"]), asql.Where("status", "=", page.FilterMap["status"]), asql.Where("is_revoke_on_refund", "=", page.FilterMap["is_revoke_on_refund"]), asql.Where("created_at", "BETWEEN", page.FilterMap["created_at"]),
		asql.Filters(page.Filter),
		asql.Order(page.Order, page.Desc),
		asql.Paginate(page.PageSize, page.CurrentPage),
	).Find(&list).Offset(-1).Limit(1).Count(&total).Error
	return list, total, err
}

func (dao *CardsDao) CountRelatedOrders(cardIds []int) (int64, error) {
	var total int64
	err := dao.db.Model(&models.Orders{}).Where("card_id IN ?", cardIds).Count(&total).Error
	return total, err
}

// GetById
func (dao *CardsDao) GetById(id int) (row models.Cards) {
	dao.db.Model(&dao.models).Where("id=?", id).Limit(1).Find(&row)
	return row
}
