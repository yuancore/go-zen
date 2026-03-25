package dto

import "github.com/yuancore/go-zen/examples/quickstart/app/entity/vo"

type PageResult[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

type UserListResult = PageResult[vo.User]
