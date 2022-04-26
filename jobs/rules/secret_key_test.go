package rules_test

import (
	"io/ioutil"
	"testing"

	"github.com/sonda2208/grc/jobs/rules"

	"github.com/sonda2208/grc/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretKeyRule(t *testing.T) {
	testCases := []struct {
		filePath      string
		expectedLines []int
	}{
		{
			filePath:      "misc/multiple_matches.go",
			expectedLines: []int{14, 27, 41, 52},
		},
		{
			filePath:      "misc/private_key_match.dart",
			expectedLines: []int{52},
		},
		{
			filePath:      "misc/no_match.go",
			expectedLines: []int{},
		},
	}

	skr, err := rules.NewSecretKeyRule()
	require.NoError(t, err)

	for _, tc := range testCases {
		filePath := util.FindFile(tc.filePath)
		fileBuf, err := ioutil.ReadFile(filePath)
		require.NoError(t, err)

		res, err := skr.Scan(filePath, fileBuf)
		require.NoError(t, err)

		var lines []int
		for _, r := range res {
			lines = append(lines, r.Location.Positions.Begin.Line)
		}

		assert.ElementsMatch(t, tc.expectedLines, lines)
	}
}
