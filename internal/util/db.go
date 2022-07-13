package util

import (
	"context"
	"database/sql"
)

func ExecWithTx(ctx context.Context, tx *sql.Tx, f func(context.Context, *sql.Tx) error) error {
	var err error
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	err = f(ctx, tx)
	return err
}
