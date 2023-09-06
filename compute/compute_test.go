package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/v1tbrah/kvdb/model"
)

func TestParseStateMachine_ParseData(t *testing.T) {
	tests := []struct {
		name           string
		inputData      string
		expectedResult model.QueryMetaData
		wantErr        bool
	}{
		{
			name:           "empty input data",
			inputData:      "",
			expectedResult: model.QueryMetaData{},
			wantErr:        true,
		},
		{
			name:           "op set. empty key",
			inputData:      "SET ",
			expectedResult: model.QueryMetaData{},
			wantErr:        true,
		},
		{
			name:           "op set. empty value",
			inputData:      "SET FOO",
			expectedResult: model.QueryMetaData{},
			wantErr:        true,
		},
		{
			name:      "op set. OK",
			inputData: "SET FOO BAR",
			expectedResult: model.QueryMetaData{
				OpType: model.OpTypeSet,
				Key:    "FOO",
				Value:  "BAR",
			},
			wantErr: false,
		},
		{
			name:           "op get. empty key",
			inputData:      "GET ",
			expectedResult: model.QueryMetaData{},
			wantErr:        true,
		},
		{
			name:      "op get. OK",
			inputData: "GET FOO",
			expectedResult: model.QueryMetaData{
				OpType: model.OpTypeGet,
				Key:    "FOO",
				Value:  "",
			},
			wantErr: false,
		},
		{
			name:      "op get. Key separated by space",
			inputData: "GET FOO BAR",
			expectedResult: model.QueryMetaData{
				OpType: model.OpTypeGet,
				Key:    "FOO BAR",
				Value:  "",
			},
			wantErr: false,
		},
		{
			name:           "op delete. empty key",
			inputData:      "DELETE ",
			expectedResult: model.QueryMetaData{},
			wantErr:        true,
		},
		{
			name:      "op delete. OK",
			inputData: "DELETE FOO",
			expectedResult: model.QueryMetaData{
				OpType: model.OpTypeDelete,
				Key:    "FOO",
				Value:  "",
			},
			wantErr: false,
		},
		{
			name:      "op delete. camel cas op type",
			inputData: "dElEtE FOO",
			expectedResult: model.QueryMetaData{
				OpType: model.OpTypeDelete,
				Key:    "FOO",
				Value:  "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			psm := NewParseStateMachine(tt.inputData)
			got, err := psm.ParseData()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, got)
			}
		})
	}
}
