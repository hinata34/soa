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
	return &model.User{}, nil
}

func (p *ProfileRepoPostgresql) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	err := p.db.Get(ctx, &user, "SELECT id, login, password, name, surname, birthday, email, mobile_number FROM users WHERE login=$1", login)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *ProfileRepoPostgresql) Update(ctx context.Context, user *model.User) error {
	var err error
	if user.Password == "" {
		_, err = p.db.Exec(ctx,
			"UPDATE users SET name=$1, surname=$2, mobile_number=$3 WHERE login=$4",
			user.Name, user.Surname, user.MobileNumber, user.Login)
	} else {
		_, err = p.db.Exec(ctx,
			"UPDATE users SET name=$1, surname=$2, mobile_number=$3, password=$4 WHERE login=$5",
			user.Name, user.Surname, user.MobileNumber, user.Password, user.Login)
	}
	return err
}

func (p *ProfileRepoPostgresql) Remove(ctx context.Context, user *model.User) error {
	return nil
}
