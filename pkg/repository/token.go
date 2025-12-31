package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/s-404/go-auth-example/pkg/entity"
)

type ITokenRepository interface {
	Create(data entity.Token) (*entity.Token, error)
	FindByUserGuid(userGuid string) (*entity.Token, error)
	FindByToken(token string) (*entity.Token, error)
	FindByGuid(guid string) (*entity.Token, error)
	Update(guid string, token string) error
	DestroyByToken(token string) error
}

type TokenRepository struct {
	db    *sqlx.DB
	table string
}

func NewTokenRepository(db *sqlx.DB, table string) *TokenRepository {
	return &TokenRepository{
		db:    db,
		table: table,
	}
}

func (r *TokenRepository) Create(data entity.Token) (*entity.Token, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_guid, token) values ($1, $2) RETURNING *", r.table)
	row := r.db.QueryRow(query, data.UserGuid, data.Token)

	var token entity.Token
	if err := row.Scan(
		&token.Guid,
		&token.Token,
		&token.UserGuid,
		&token.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) FindByUserGuid(userGuid string) (*entity.Token, error) {
	var token entity.Token
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_guid=$1", r.table)
	err := r.db.Get(&token, query, userGuid)

	return &token, err
}

func (r *TokenRepository) FindByToken(tokenValue string) (*entity.Token, error) {
	var token entity.Token
	query := fmt.Sprintf("SELECT * FROM %s WHERE token=$1", r.table)
	err := r.db.Get(&token, query, tokenValue)

	return &token, err
}

func (r *TokenRepository) FindByGuid(guid string) (*entity.Token, error) {
	var token entity.Token
	query := fmt.Sprintf("SELECT * FROM %s WHERE guid=$1", r.table)
	err := r.db.Get(&token, query, guid)

	return &token, err
}

func (r *TokenRepository) Update(guid string, token string) error {
	_, err := r.db.NamedExec(fmt.Sprintf(`UPDATE %s SET token=:token WHERE guid=:guid`, r.table),
		map[string]interface{}{
			"token": token,
			"guid":  guid,
		})

	return err
}

func (r *TokenRepository) DestroyByToken(token string) error {
	_, err := r.db.NamedExec(fmt.Sprintf(`DELETE FROM %s WHERE token=:token`, r.table),
		map[string]interface{}{
			"token": token,
		})

	return err
}
