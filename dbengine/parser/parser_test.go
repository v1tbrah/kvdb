package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompute(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "empty string",
			in:   "",
			want: []string{},
		},
		{
			name: "string only with spaces",
			in:   "   ",
			want: []string{},
		},
		{
			name: "one word",
			in:   "SET",
			want: []string{"SET"},
		},
		{
			name: "one word with spaces in start",
			in:   "   SET",
			want: []string{"SET"},
		},
		{
			name: "one word with spaces in end",
			in:   "SET   ",
			want: []string{"SET"},
		},
		{
			name: "one word with spaces in start and end",
			in:   "   SET   ",
			want: []string{"SET"},
		},
		{
			name: "string with spaces in start",
			in:   "   SET FOO BAR",
			want: []string{"SET", "FOO", "BAR"},
		},
		{
			name: "string with spaces in end",
			in:   "SET FOO BAR   ",
			want: []string{"SET", "FOO", "BAR"},
		},
		{
			name: "string with spaces in start and end",
			in:   "   SET FOO BAR   ",
			want: []string{"SET", "FOO", "BAR"},
		},
		{
			name: "simple string",
			in:   "SET FOO BAR",
			want: []string{"SET", "FOO", "BAR"},
		},
		{
			name: "string with a lot spaces between words",
			in:   "SET   FOO  BAR",
			want: []string{"SET", "FOO", "BAR"},
		},
		// TODO add cases with runes which contain more than one byte
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Compute(tt.in))
		})
	}
}
