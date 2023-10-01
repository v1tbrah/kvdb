package dbengine

import (
	"errors"
	"fmt"

	"github.com/v1tbrah/kvdb/dbengine/parser"
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

func (e *Engine) Process(in string) (out string, err error) {
	if len(in) == 0 {
		return out, errors.New("empty input")
	}

	parsed := parser.Compute(in)
	if len(parsed) == 0 {
		return out, errors.New("empty parsed input")
	}

	executableOperation, err := e.analyzeOperation(parsed)
	if err != nil {
		return out, fmt.Errorf("e.analyzeOperation: %w", err)
	}

	out = executableOperation()

	return out, nil
}
