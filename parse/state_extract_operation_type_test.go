package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/v1tbrah/kvdb/model"
)

func Test_extractOperationTypeState_toggle(t *testing.T) {
	tErrState := errorState{}

	tests := []struct {
		name string
		psm  *parseStateMachine

		expectedCurrentState parseState
		expectedOpType       model.OpType
	}{
		{
			name:                 "invalid op type",
			psm:                  newParseStateMachine("SED"),
			expectedCurrentState: tErrState,
		},
		{
			name: "SET op type",
			psm: func() *parseStateMachine {
				psm := newParseStateMachine("SET")
				psm.operationType = "SET"
				return psm
			}(),
			expectedOpType:       model.OpTypeSet,
			expectedCurrentState: extractKeyState{},
		},
		{
			name: "GET op type",
			psm: func() *parseStateMachine {
				psm := newParseStateMachine("GET")
				psm.operationType = "GET"
				return psm
			}(),
			expectedOpType:       model.OpTypeGet,
			expectedCurrentState: extractKeyState{},
		},
		{
			name: "DELETE op type",
			psm: func() *parseStateMachine {
				psm := newParseStateMachine("DELETE")
				psm.operationType = "DELETE"
				return psm
			}(),
			expectedOpType:       model.OpTypeDelete,
			expectedCurrentState: extractKeyState{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eot := extractOperationTypeState{}
			eot.toggle(tt.psm)
			assert.Equal(t, tt.expectedCurrentState, tt.psm.currentState)
			if tt.expectedCurrentState != tErrState {
				assert.Equal(t, tt.expectedOpType, tt.psm.parsedData.OpType)
			}
		})
	}
}
