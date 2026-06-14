package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upFixBitDepthNulls, downFixBitDepthNulls)
}

func upFixBitDepthNulls(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "UPDATE media_file SET bit_depth = 0 WHERE bit_depth IS NULL")
	return err
}

func downFixBitDepthNulls(ctx context.Context, tx *sql.Tx) error {
	return nil
}
