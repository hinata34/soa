package postgresql

import (
	"context"
	"promo/internal/userservice/db"
	model "promo/internal/userservice/models"
)

type ProfileRepoPostgresql struct {
	db db.DBops
}

func NewProfileRepoPostgresql(db db.DBops) *ProfileRepoPostgresql {
	return &ProfileRepoPostgresql{db: db}
}

func (p *ProfileRepoPostgresql) GetById(ctx context.Context, id uint64) (*model.User, error) {

}

func (p *ProfileRepoPostgresql) GetByLogin(ctx context.Context, login string) (*model.User, error) {

}

func (p *ProfileRepoPostgresql) Update(ctx context.Context, user *model.User) error {

}

func (p *ProfileRepoPostgresql) Remove(ctx context.Context, user *model.User) error {

}
