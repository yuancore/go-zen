package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"

	zdb "github.com/yuancore/go-zen/adapter/db/gorm"
	"github.com/yuancore/go-zen/examples/quickstart/app/entity/models"
	"github.com/yuancore/go-zen/zen"
	"gorm.io/gorm"
)

type UserListFilter struct {
	Page    int
	Size    int
	Keyword string
	Status  *int8
}

type UserDAO struct {
	app *zen.App
}

func NewUserDAO(app *zen.App) *UserDAO {
	return &UserDAO{app: app}
}

func (d *UserDAO) AutoMigrate(ctx context.Context) error {
	return d.db(ctx).AutoMigrate(&models.User{})
}

func (d *UserDAO) Create(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("create user: nil entity")
	}
	return d.db(ctx).Create(user).Error
}

func (d *UserDAO) Update(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("update user: nil entity")
	}
	tx := d.db(ctx).Save(user)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (d *UserDAO) Delete(ctx context.Context, id uint64) error {
	tx := d.db(ctx).Delete(&models.User{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (d *UserDAO) GetByID(ctx context.Context, id uint64) (*models.User, error) {
	var user models.User
	if err := d.db(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDAO) List(ctx context.Context, filter UserListFilter) ([]models.User, int64, error) {
	query := d.db(ctx).Model(&models.User{})

	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", pattern, pattern, pattern)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	if total == 0 {
		return []models.User{}, 0, nil
	}

	users := make([]models.User, 0, filter.Size)
	if err := query.
		Order("id DESC").
		Offset((filter.Page - 1) * filter.Size).
		Limit(filter.Size).
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	return users, total, nil
}

func (d *UserDAO) db(ctx context.Context) *gorm.DB {
	db := zdb.MustResolve(d.app)
	if ctx == nil {
		return db
	}
	return db.WithContext(ctx)
}
