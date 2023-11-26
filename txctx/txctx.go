package txctx

import (
	"context"

	"github.com/google/uuid"
)

type txCtxKey struct{}

func CtxWithTx(ctx context.Context) context.Context {
	return context.WithValue(ctx, txCtxKey{}, uuid.NewString())
}

func Tx(ctx context.Context) (tx string) {
	tx, _ = ctx.Value(txCtxKey{}).(string)
	return tx
}
