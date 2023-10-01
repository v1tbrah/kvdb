package dbengine

import (
	"context"
	"errors"
	"log/slog"

	"github.com/v1tbrah/kvdb/dbengine/parser"
	"github.com/v1tbrah/kvdb/txctx"
)

type memory interface {
	Set(key, value string)
	Get(key string) string
	Delete(key string)
}

type Engine struct {
	memory memory
}

func New(memory memory) (*Engine, error) {
	if memory == nil {
		return nil, errors.New("memory is nil")
	}

	return &Engine{
		memory: memory,
	}, nil
}

func (e *Engine) Process(ctx context.Context, in string) (out string, err error) {
	if len(in) == 0 {
		slog.WarnContext(ctx, "empty input",
			slog.String("tx", txctx.Tx(ctx)))
		return out, errors.New("invalid input")
	}

	parsed := parser.Compute(in)
	if len(parsed) == 0 {
		slog.WarnContext(ctx, "empty parsed input",
			slog.String("tx", txctx.Tx(ctx)), slog.String("input", in))
		return out, errors.New("invalid input")
	}

	executableOperationFn, err := e.analyzeOperation(parsed)
	if err != nil {
		slog.WarnContext(ctx, "analyzeOperation",
			slog.String("tx", txctx.Tx(ctx)), slog.Any("parsed", parsed),
			slog.String("error", err.Error()))
		return out, errors.New("invalid input")
	}

	out = executableOperationFn()

	return out, nil
}
