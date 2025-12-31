package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/s-404/go-auth-example/pkg/entity"
)

type IUserRepository interface {
	Create(data entity.User) (*entity.User, error)
	FindByGuid(guid string) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	FindByUsernameAndPassword(username string, password string) (*entity.User, error)
}

type UserRepository struct {
	db    *sqlx.DB
	table string
}

func NewUserRepository(db *sqlx.DB, table string) *UserRepository {
	return &UserRepository{
		db:    db,
		table: table,
	}
}

func (r *UserRepository) Create(data entity.User) (*entity.User, error) {
	query := fmt.Sprintf("INSERT INTO %s (username, password) values ($1, $2) RETURNING guid, username", r.table)
	row := r.db.QueryRow(query, data.Username, data.Password)

	var user entity.User
	if err := row.Scan(&user.Guid, &user.Username); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByGuid(guid string) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE guid=$1", r.table)
	err := r.db.Get(&user, query, guid)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1", r.table)
	err := r.db.Get(&user, query, username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByUsernameAndPassword(username string, password string) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1 AND password=$2", r.table)
	err := r.db.Get(&user, query, username, password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
