package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateTableTokens, downCreateTableTokens)
}

func upCreateTableTokens(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE tokens
		(
			guid uuid NOT NULL DEFAULT gen_random_uuid(),
			user_guid uuid NOT NULL
				REFERENCES users (guid)
				ON DELETE CASCADE
				ON UPDATE CASCADE,
			token TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(user_guid, token),
		
			CONSTRAINT tokens_pkey PRIMARY KEY (guid)
		);`)
	if err != nil {
		return fmt.Errorf("failed create table tokens: %w", err)
	}

	return nil
}

func downCreateTableTokens(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx,
		"DROP TABLE IF EXISTS tokens;")
	if err != nil {
		return fmt.Errorf("failed drop table tokens: %w", err)
	}

	return nil
}
