package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateTableUsers, downCreateTableUsers)
}

func upCreateTableUsers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE users
		(
			guid uuid NOT NULL DEFAULT gen_random_uuid(),
			username varchar(20) NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			UNIQUE (username),
		
			CONSTRAINT users_pkey PRIMARY KEY (guid)
		);`)
	if err != nil {
		return fmt.Errorf("failed create table users: %w", err)
	}

	return nil
}

func downCreateTableUsers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx,
		"DROP TABLE IF EXISTS users;")
	if err != nil {
		return fmt.Errorf("failed drop table users: %w", err)
	}

	return nil
}
