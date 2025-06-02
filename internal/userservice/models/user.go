package model

import "time"

type User struct {
	ID           uint64    `db:"id"`
	Login        string    `db:"login"`
	Password     string    `db:"password"`
	Name         string    `db:"name"`
	Surname      string    `db:"surname"`
	Birthday     time.Time `db:"birthday"`
	Email        string    `db:"email"`
	MobileNumber string    `db:"mobile_number"`
	Created      time.Time `db:"created"`
	Updated      time.Time `db:"updated"`
}
