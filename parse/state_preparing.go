package parse

import (
	"errors"
	"strings"

	"github.com/v1tbrah/kvdb/model"
)

type preparingState struct{}

func (v preparingState) toggle(psm *parseStateMachine) {
	components := strings.SplitN(psm.dataToParse, " ", 2)
	if len(components) < 2 {
		psm.parsingError = errors.New("query must contain operation and key, separated by space ' '")
		psm.currentState = errorState{}
		return
	}

	psm.operationType = strings.ToUpper(components[0])
	if psm.operationType != string(model.OpTypeSet) {
		psm.key = components[1]
	} else {
		setComponents := strings.SplitN(components[1], " ", 2)
		psm.key = setComponents[0]
		if len(setComponents) > 1 {
			psm.value = setComponents[1]
		}
	}

	psm.currentState = extractOperationTypeState{}
}
