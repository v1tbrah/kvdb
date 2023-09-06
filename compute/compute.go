package compute

import (
	"github.com/v1tbrah/kvdb/model"
)

type parseState interface {
	toggle(psm *ParseStateMachine)
}

type ParseStateMachine struct {
	currentState parseState

	dataToParse string

	operationType string
	key           string
	value         string

	parsedData model.QueryMetaData

	parsingError error
}

func NewParseStateMachine(dataToParse string) *ParseStateMachine {
	return &ParseStateMachine{
		currentState: preparingState{},
		dataToParse:  dataToParse,
		parsedData:   model.QueryMetaData{},
		parsingError: nil,
	}
}

func (psm *ParseStateMachine) ParseData() (model.QueryMetaData, error) {
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

type parsedState struct{}

func (p parsedState) toggle(_ *ParseStateMachine) {
	panic("toggle can't be called after parsed state. check for it and process")
}

type errorState struct{}

func (e errorState) toggle(_ *ParseStateMachine) {
	panic("toggle can't be called after error state. check for it and process")
}
