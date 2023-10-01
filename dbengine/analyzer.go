package dbengine

import (
	"errors"
	"fmt"
	"strings"
)

type OpType string

const (
	OpTypeSet    OpType = "SET"
	OpTypeGet    OpType = "GET"
	OpTypeDelete OpType = "DELETE"
)

func (ot OpType) isValid() bool {
	return ot == OpTypeSet ||
		ot == OpTypeGet ||
		ot == OpTypeDelete
}

type executableOperation func() (response string)

func (e *Engine) analyzeOperation(in []string) (operationFn executableOperation, err error) {
	if len(in) == 0 {
		return operationFn, errors.New("empty input")
	}

	opType := OpType(strings.ToUpper(in[0]))
	if !opType.isValid() {
		return operationFn, errors.New("invalid operation type")
	}

	arguments := in[1:]
	if len(arguments) == 0 {
		return operationFn, errors.New("empty arguments")
	}

	switch opType {
	case OpTypeGet:
		if len(arguments) != 1 {
			return operationFn, fmt.Errorf("unexpected count of arguments, expeсted 1, got %d", len(arguments))
		}
		operationFn = func() string {
			return e.operationGet(arguments[0])
		}
	case OpTypeSet:
		if len(arguments) != 2 {
			return operationFn, fmt.Errorf("unexpected count of arguments, expeсted 2, got %d", len(arguments))
		}
		operationFn = func() string {
			return e.operationSet(arguments[0], arguments[1])
		}
	case OpTypeDelete:
		if len(arguments) != 1 {
			return operationFn, fmt.Errorf("unexpected count of arguments, expeсted 1, got %d", len(arguments))
		}
		operationFn = func() string {
			return e.operationDelete(arguments[0])
		}
	default:
		err = errors.New("unsupported operation type")
	}

	return operationFn, nil
}

const responseSuccessful = "OK"

func (e *Engine) operationGet(key string) (response string) {
	response = e.memory.Get(key)

	return response
}

func (e *Engine) operationSet(key, value string) (response string) {
	e.memory.Set(key, value)

	return responseSuccessful
}

func (e *Engine) operationDelete(key string) (response string) {
	e.memory.Delete(key)

	return responseSuccessful
}
