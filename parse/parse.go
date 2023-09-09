package parse

import (
	"github.com/v1tbrah/kvdb/model"
)

func ParseData(dataToParse string) (model.QueryMetaData, error) {
	psm := newParseStateMachine(dataToParse)

	var eState errorState
	var pState parsedState
	for {
		psm.currentState.toggle(psm)
		if psm.currentState == eState {
			return model.QueryMetaData{}, psm.parsingError
		}
		if psm.currentState == pState {
			return psm.parsedData, nil
		}
	}
}

type parseState interface {
	toggle(psm *parseStateMachine)
}

type parseStateMachine struct {
	currentState parseState

	dataToParse string

	operationType string
	key           string
	value         string

	parsedData model.QueryMetaData

	parsingError error
}

func newParseStateMachine(dataToParse string) *parseStateMachine {
	return &parseStateMachine{
		currentState: preparingState{},
		dataToParse:  dataToParse,
		parsedData:   model.QueryMetaData{},
		parsingError: nil,
	}
}

type parsedState struct{}

func (p parsedState) toggle(_ *parseStateMachine) {
	panic("toggle can't be called after parsed state. check for it and process")
}

type errorState struct{}

func (e errorState) toggle(_ *parseStateMachine) {
	panic("toggle can't be called after error state. check for it and process")
}
