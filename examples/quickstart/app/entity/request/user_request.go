package request

import (
	"strings"

	quickentity "github.com/yuancore/go-zen/examples/quickstart/app/entity"
)

type ListUsersRequest struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	Size    int    `form:"size" binding:"omitempty,min=1,max=100"`
	Keyword string `form:"keyword" binding:"omitempty,max=64"`
	Status  *int8  `form:"status" binding:"omitempty,oneof=0 1"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Nickname string `json:"nickname" binding:"omitempty,max=64"`
	Email    string `json:"email" binding:"required,email,max=128"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Password string `json:"password" binding:"required,min=12,max=72"`
	Status   *int8  `json:"status" binding:"omitempty,oneof=0 1"`
}

type UpdateUserRequest struct {
	Nickname *string `json:"nickname" binding:"omitempty,max=64"`
	Email    *string `json:"email" binding:"omitempty,email,max=128"`
	Phone    *string `json:"phone" binding:"omitempty,max=20"`
	Password *string `json:"password" binding:"omitempty,min=12,max=72"`
	Status   *int8   `json:"status" binding:"omitempty,oneof=0 1"`
}

func (r *ListUsersRequest) Normalize() {
	r.Page, r.Size = quickentity.NormalizePage(r.Page, r.Size)
	r.Keyword = strings.TrimSpace(r.Keyword)
}

func (r UpdateUserRequest) HasChanges() bool {
	return r.Nickname != nil || r.Email != nil || r.Phone != nil || r.Password != nil || r.Status != nil
}
