package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/v1tbrah/kvdb/model"
)

func Test_extractKeyState_toggle(t *testing.T) {
	tErrState := errorState{}

	tests := []struct {
		name string
		psm  *parseStateMachine

		expectedCurrentState parseState
		expectedKey          string
	}{
		{
			name:                 "empty key",
			psm:                  newParseStateMachine("SET"),
			expectedCurrentState: tErrState,
		},
		{
			name: "op set. valid key",
			psm: func() *parseStateMachine {
				psm := newParseStateMachine("SET FOO BAR")
				psm.operationType = "SET"
				psm.parsedData.OpType = model.OpTypeSet
				psm.key = "FOO"
				return psm
			}(),
			expectedKey:          "FOO",
			expectedCurrentState: extractValueState{},
		},
		{
			name: "op get. valid key",
			psm: func() *parseStateMachine {
				psm := newParseStateMachine("GET FOO")
				psm.operationType = "GET"
				psm.parsedData.OpType = model.OpTypeGet
				psm.key = "FOO"
				return psm
			}(),
			expectedKey:          "FOO",
			expectedCurrentState: parsedState{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eks := extractKeyState{}
			eks.toggle(tt.psm)
			assert.Equal(t, tt.expectedCurrentState, tt.psm.currentState)
			if tt.expectedCurrentState != tErrState {
				assert.Equal(t, tt.expectedKey, tt.psm.parsedData.Key)
			}
		})
	}
}
