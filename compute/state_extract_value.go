package compute

import "errors"

type extractValueState struct{}

func (ev extractValueState) toggle(psm *parseStateMachine) {
	if psm.value == "" {
		psm.parsingError = errors.New("value must not be empty")
		psm.currentState = errorState{}
		return
	}

	psm.parsedData.Value = psm.value

	psm.currentState = parsedState{}
}
