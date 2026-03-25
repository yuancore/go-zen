package api

import (
	"errors"
	"net/http"

	"github.com/yuancore/go-zen/examples/quickstart/app/entity/request"
	"github.com/yuancore/go-zen/examples/quickstart/app/service"
	"github.com/yuancore/go-zen/zen"
)

type UserAPI struct {
	service *service.UserService
	logger  zen.Logger
}

func NewUserAPI(service *service.UserService, logger zen.Logger) *UserAPI {
	return &UserAPI{
		service: service,
		logger:  logger.With("module", "user_api"),
	}
}

func (a *UserAPI) List(c zen.Context) {
	var req request.ListUsersRequest
	if !bindQuery(c, &req) {
		return
	}

	result, err := a.service.List(c.Request().Context(), req)
	if err != nil {
		a.handleError(c, "list", err)
		return
	}
	writeOK(c, result)
}

func (a *UserAPI) Get(c zen.Context) {
	id, ok := parseUint64ID(c, "id")
	if !ok {
		return
	}

	user, err := a.service.Get(c.Request().Context(), id)
	if err != nil {
		a.handleError(c, "get", err)
		return
	}
	writeOK(c, user)
}

func (a *UserAPI) Create(c zen.Context) {
	var req request.CreateUserRequest
	if !bindJSON(c, &req) {
		return
	}

	user, err := a.service.Create(c.Request().Context(), req)
	if err != nil {
		a.handleError(c, "create", err)
		return
	}
	writeCreated(c, user)
}

func (a *UserAPI) Update(c zen.Context) {
	id, ok := parseUint64ID(c, "id")
	if !ok {
		return
	}

	var req request.UpdateUserRequest
	if !bindJSON(c, &req) {
		return
	}

	user, err := a.service.Update(c.Request().Context(), id, req)
	if err != nil {
		a.handleError(c, "update", err)
		return
	}
	writeOK(c, user)
}

func (a *UserAPI) Delete(c zen.Context) {
	id, ok := parseUint64ID(c, "id")
	if !ok {
		return
	}

	if err := a.service.Delete(c.Request().Context(), id); err != nil {
		a.handleError(c, "delete", err)
		return
	}
	writeOK(c, map[string]bool{"deleted": true})
}

func (a *UserAPI) handleError(c zen.Context, operation string, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		writeError(c, http.StatusNotFound, "user not found")
	case errors.Is(err, service.ErrUserConflict):
		writeError(c, http.StatusConflict, "user already exists")
	case errors.Is(err, service.ErrInvalidUserData):
		writeError(c, http.StatusBadRequest, "invalid user data")
	default:
		a.logger.Error("user request failed", "op", operation, "err", err)
		writeError(c, http.StatusInternalServerError, "internal server error")
	}
}
