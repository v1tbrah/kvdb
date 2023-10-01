package parser

import (
	"strings"
)

const space = ' '

func Compute(dataToParse string) []string {
	result := make([]string, 0)

	preparedData := strings.TrimSpace(dataToParse)

	var (
		builder     strings.Builder
		isPrevSpace bool
	)

	for _, v := range preparedData {
		if v == space {
			if !isPrevSpace {
				result = append(result, builder.String())
				builder.Reset()
			}

			isPrevSpace = true
			continue
		}

		isPrevSpace = false

		_, _ = builder.WriteRune(v) // the error can never come back. see source code
	}

	lastWord := builder.String()
	if lastWord != "" {
		result = append(result, lastWord)
	}

	return result
}
