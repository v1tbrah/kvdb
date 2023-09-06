package compute

import (
	"errors"

	"github.com/v1tbrah/kvdb/model"
)

type extractOperationTypeState struct{}

func (eot extractOperationTypeState) toggle(psm *ParseStateMachine) {
	switch model.OpType(psm.operationType) {
	case model.OpTypeSet:
		psm.parsedData.OpType = model.OpTypeSet
	case model.OpTypeGet:
		psm.parsedData.OpType = model.OpTypeGet
	case model.OpTypeDelete:
		psm.parsedData.OpType = model.OpTypeDelete
	default:
		psm.parsingError = errors.New("invalid operation. supported only 'SET', 'GET', 'DELETE'")
		psm.currentState = errorState{}
		return
	}

	psm.currentState = extractKeyState{}
}
