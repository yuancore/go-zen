package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/yuancore/go-zen/examples/quickstart/app/dao"
	"github.com/yuancore/go-zen/examples/quickstart/app/entity/dto"
	"github.com/yuancore/go-zen/examples/quickstart/app/entity/models"
	"github.com/yuancore/go-zen/examples/quickstart/app/entity/request"
	"github.com/yuancore/go-zen/examples/quickstart/app/entity/vo"
	"github.com/yuancore/go-zen/zen"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserConflict    = errors.New("user conflict")
	ErrInvalidUserData = errors.New("invalid user data")

	usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]{3,32}$`)
)

type UserService struct {
	dao    *dao.UserDAO
	logger zen.Logger
}

func NewUserService(dao *dao.UserDAO, logger zen.Logger) *UserService {
	return &UserService{
		dao:    dao,
		logger: logger.With("module", "user_service"),
	}
}

func (s *UserService) Migrate(ctx context.Context) error {
	if err := s.dao.AutoMigrate(ctx); err != nil {
		return fmt.Errorf("auto migrate users: %w", err)
	}
	return nil
}

func (s *UserService) List(ctx context.Context, req request.ListUsersRequest) (dto.UserListResult, error) {
	req.Normalize()

	users, total, err := s.dao.List(ctx, dao.UserListFilter{
		Page:    req.Page,
		Size:    req.Size,
		Keyword: req.Keyword,
		Status:  req.Status,
	})
	if err != nil {
		return dto.UserListResult{}, err
	}

	items := make([]vo.User, 0, len(users))
	for _, user := range users {
		items = append(items, vo.NewUser(user))
	}

	return dto.UserListResult{
		Items: items,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

func (s *UserService) Get(ctx context.Context, id uint64) (vo.User, error) {
	user, err := s.dao.GetByID(ctx, id)
	if err != nil {
		return vo.User{}, s.mapError("get user", err)
	}
	return vo.NewUser(*user), nil
}

func (s *UserService) Create(ctx context.Context, req request.CreateUserRequest) (vo.User, error) {
	username := strings.TrimSpace(req.Username)
	if !usernamePattern.MatchString(username) {
		return vo.User{}, ErrInvalidUserData
	}

	email := normalizeEmail(req.Email)
	if email == "" {
		return vo.User{}, ErrInvalidUserData
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return vo.User{}, err
	}

	user := models.User{
		Username:     username,
		Nickname:     strings.TrimSpace(req.Nickname),
		Email:        email,
		Phone:        strings.TrimSpace(req.Phone),
		PasswordHash: passwordHash,
		Status:       defaultStatus(req.Status),
	}

	if err := s.dao.Create(ctx, &user); err != nil {
		return vo.User{}, s.mapError("create user", err)
	}

	return vo.NewUser(user), nil
}

func (s *UserService) Update(ctx context.Context, id uint64, req request.UpdateUserRequest) (vo.User, error) {
	if !req.HasChanges() {
		return vo.User{}, ErrInvalidUserData
	}

	user, err := s.dao.GetByID(ctx, id)
	if err != nil {
		return vo.User{}, s.mapError("load user", err)
	}

	if req.Nickname != nil {
		user.Nickname = strings.TrimSpace(*req.Nickname)
	}
	if req.Email != nil {
		email := normalizeEmail(*req.Email)
		if email == "" {
			return vo.User{}, ErrInvalidUserData
		}
		user.Email = email
	}
	if req.Phone != nil {
		user.Phone = strings.TrimSpace(*req.Phone)
	}
	if req.Password != nil {
		passwordHash, err := hashPassword(*req.Password)
		if err != nil {
			return vo.User{}, err
		}
		user.PasswordHash = passwordHash
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.dao.Update(ctx, user); err != nil {
		return vo.User{}, s.mapError("update user", err)
	}

	return vo.NewUser(*user), nil
}

func (s *UserService) Delete(ctx context.Context, id uint64) error {
	return s.mapError("delete user", s.dao.Delete(ctx, id))
}

func (s *UserService) mapError(operation string, err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrUserNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrUserConflict
	case errors.Is(err, bcrypt.ErrPasswordTooLong):
		return ErrInvalidUserData
	default:
		s.logger.Error("user service failed", "op", operation, "err", err)
		return fmt.Errorf("%s: %w", operation, err)
	}
}

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func hashPassword(password string) (string, error) {
	password = strings.TrimSpace(password)
	if len(password) < 12 || len(password) > 72 {
		return "", ErrInvalidUserData
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hashed), nil
}

func defaultStatus(status *int8) int8 {
	if status == nil {
		return 1
	}
	return *status
}
