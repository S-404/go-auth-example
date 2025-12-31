package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	User  IUserRepository
	Token ITokenRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:  NewUserRepository(db, "users"),
		Token: NewTokenRepository(db, "tokens"),
	}
}
