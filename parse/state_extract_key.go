package parse

import (
	"errors"
	"strings"

	"github.com/v1tbrah/kvdb/model"
)

type extractKeyState struct{}

func (ek extractKeyState) toggle(psm *parseStateMachine) {
	if strings.TrimSpace(psm.key) == "" {
		psm.currentState = errorState{}
		psm.parsingError = errors.New("key must have one or more non space symbol")
		return
	}

	psm.parsedData.Key = psm.key

	if psm.parsedData.OpType == model.OpTypeSet {
		psm.currentState = extractValueState{}
		return
	}

	psm.currentState = parsedState{}
}
