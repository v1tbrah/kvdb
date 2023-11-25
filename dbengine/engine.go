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

type wal interface {
	Save(ctx context.Context, in string) error
}

type Engine struct {
	memory memory
	wal    wal
}

func New(memory memory, wal wal) (*Engine, error) {
	if memory == nil {
		return nil, errors.New("memory is nil")
	}
	if wal == nil {
		return nil, errors.New("wal is nil")
	}

	return &Engine{
		memory: memory,
		wal:    wal,
	}, nil
}

func (e *Engine) Process(ctx context.Context, in string) (out string, err error) {
	if len(in) == 0 {
		slog.WarnContext(ctx, "empty input",
			slog.String("tx", txctx.Tx(ctx)))
		return out, errors.New("invalid input")
	}

	// parse input tokens
	parsed := parser.Compute(in)
	if len(parsed) == 0 {
		slog.WarnContext(ctx, "empty parsed input",
			slog.String("tx", txctx.Tx(ctx)), slog.String("input", in))
		return out, errors.New("invalid input")
	}

	// analyze executable operation and operands
	executableOperationFn, opType, err := e.analyzeOperation(parsed)
	if err != nil {
		slog.WarnContext(ctx, "analyzeOperation",
			slog.String("tx", txctx.Tx(ctx)), slog.Any("parsed", parsed),
			slog.String("error", err.Error()))
		return out, errors.New("invalid input")
	}

	// save only write operations to WAL (write ahead log)
	if opType == OpTypeSet || opType == OpTypeDelete {
		if err = e.wal.Save(ctx, in); err != nil {
			slog.WarnContext(ctx, "wal.Save",
				slog.String("tx", txctx.Tx(ctx)), slog.String("in", in),
				slog.String("error", err.Error()))
			return out, errors.New("internal error") // TODO
		}
	}

	// execute in-mem operation
	out = executableOperationFn()

	return out, nil
}
