package util

import (
	"context"

	"github.com/duongcongtoai/toytoytoy/internal/common"
)

func ExecWithTx(ctx context.Context, tx common.Tx, f func(context.Context, common.Tx) error) error {
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
