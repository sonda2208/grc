package scanners

import (
	"bytes"
	"io/ioutil"
	"regexp"

	"github.com/sonda2208/guardrails-challenge/model"
)

func RegexSearch(re *regexp.Regexp, path string) ([]*model.Finding, error) {
	var findings []*model.Finding

	fileBuf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	matchBuf := fileBuf

	locs := re.FindAllIndex(matchBuf, -1)
	lastStart := 0
	lastLineNumber := 0
	lastMatchIndex := 0
	lastLineStartIndex := 0
	for _, match := range locs {
		start := match[0]
		lineStart := lastLineStartIndex
		if idx := bytes.LastIndex(matchBuf[lastStart:start], []byte{'\n'}); idx >= 0 {
			lineStart = lastStart + idx + 1
		}

		lineNumber, matchIndex := computeLineNumber(matchBuf, lastLineNumber, lastMatchIndex, lineStart, start)

		lastMatchIndex = matchIndex
		lastLineNumber = lineNumber

		findings = append(findings, &model.Finding{
			Type: "sast",
			Location: model.FindingLocation{
				Path: path,
				Positions: model.FindingPosition{
					Begin: model.PositionIndex{
						Line: lineNumber + 1,
					},
				},
			},
		})
	}

	return findings, nil
}

func computeLineNumber(fileBuf []byte, lastLineNumber, lastMatchIndex, lineStart, matchStart int) (lineNumber, matchIndex int) {
	lineNumber = lastLineNumber + bytes.Count(fileBuf[lastMatchIndex:matchStart], []byte{'\n'})
	return lineNumber, lineStart
}
