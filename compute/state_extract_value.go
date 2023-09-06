package compute

import "errors"

type extractValueState struct{}

func (ev extractValueState) toggle(psm *ParseStateMachine) {
	if psm.value == "" {
		psm.parsingError = errors.New("value must not be empty")
		psm.currentState = errorState{}
		return
	}

	psm.parsedData.Value = psm.value

	psm.currentState = parsedState{}
}
