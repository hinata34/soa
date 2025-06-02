//go:generate mockgen -source=./repository.go -destination=./mocks/repository.go -package=mock
package auth

import (
	"context"
	"errors"
	model "promo/internal/userservice/models"
)

var (
	ErrObjectNotFound = errors.New("object not found")
)

type AuthRepo interface {
	Add(ctx context.Context, user *model.User) (uint64, error)
	GetById(ctx context.Context, id uint64) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}
