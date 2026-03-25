package vo

import (
	"time"

	"github.com/yuancore/go-zen/examples/quickstart/app/entity/models"
)

type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(model models.User) User {
	return User{
		ID:        model.ID,
		Username:  model.Username,
		Nickname:  model.Nickname,
		Email:     model.Email,
		Phone:     model.Phone,
		Status:    model.Status,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
