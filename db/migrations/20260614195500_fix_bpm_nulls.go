package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upFixBpmNulls, downFixBpmNulls)
}

func upFixBpmNulls(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "UPDATE media_file SET bpm = 0 WHERE bpm IS NULL")
	return err
}

func downFixBpmNulls(ctx context.Context, tx *sql.Tx) error {
	return nil
}
