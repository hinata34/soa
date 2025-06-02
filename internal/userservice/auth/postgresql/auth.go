package postgresql

import (
	"context"
	"promo/internal/userservice/db"
	model "promo/internal/userservice/models"
)

type AuthRepoPostgresql struct {
	db db.DBops
}

func NewAuthRepoPostgresql(db db.DBops) *AuthRepoPostgresql {
	return &AuthRepoPostgresql{db: db}
}

func (a *AuthRepoPostgresql) Add(ctx context.Context, user *model.User) (uint64, error) {
	var id uint64
	err := a.db.ExecQueryRow(ctx,
		`INSERT INTO users(login, password, name, surname, birthday, email, mobile_number) 
		VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		user.Login, user.Password, user.Name, user.Surname, user.Birthday, user.Email, user.MobileNumber,
	).Scan(&id)
	return id, err
}

func (a *AuthRepoPostgresql) GetById(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := a.db.Get(ctx, &user, "SELECT id, login, password, name, surname, birthday, email, mobile_number FROM users WHERE id=$1", id) // think about name of table
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *AuthRepoPostgresql) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	err := a.db.Get(ctx, &user, "SELECT id, login, password, name, surname, birthday, email, mobile_number FROM users WHERE login=$1", login)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *AuthRepoPostgresql) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := a.db.Get(ctx, &user, "SELECT id, login, password, name, surname, birthday, email, mobile_number FROM users WHERE email=$1", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
