package scanners_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sonda2208/guardrails-challenge/util"

	"github.com/sonda2208/guardrails-challenge/scanners"

	"github.com/sonda2208/guardrails-challenge/scanners/secret_key_scanner"

	"github.com/stretchr/testify/require"
)

func TestRegexSearchFile(t *testing.T) {
	testCases := []struct {
		filePath      string
		expectedLines []int
		isError       bool
	}{
		{
			filePath: "misc/abc.go",
			isError:  true,
		},
		{
			filePath:      "misc/multiple_matches.go",
			expectedLines: []int{14, 27, 41, 52},
			isError:       false,
		},
		{
			filePath:      "misc/private_key_match.dart",
			expectedLines: []int{52},
			isError:       false,
		},
		{
			filePath:      "misc/no_match.go",
			expectedLines: []int{},
			isError:       false,
		},
	}

	re, err := regexp.Compile(secret_key_scanner.SecretKeysPattern)
	require.NoError(t, err)

	for _, tc := range testCases {
		filePath := util.FindFile(tc.filePath)
		res, err := scanners.RegexSearch(re, filePath)
		if tc.isError {
			assert.Error(t, err)
		} else {
			require.NoError(t, err)

			var lines []int
			for _, r := range res {
				lines = append(lines, r.Location.Positions.Begin.Line)
			}

			assert.ElementsMatch(t, tc.expectedLines, lines)
		}
	}
}
