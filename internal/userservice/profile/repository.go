//go:generate mockgen -source=./repository.go -destination=./mocks/repository.go -package=mock_repository
package profile

import (
	"context"
	model "promo/internal/userservice/models"
)

type ProfileRepo interface {
	GetById(ctx context.Context, id uint64) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error) // profile acquiring
	Update(ctx context.Context, user *model.User) error
	Remove(ctx context.Context, user *model.User) error
}
