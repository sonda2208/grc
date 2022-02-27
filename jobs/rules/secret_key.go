package rules

import (
	"bytes"
	"regexp"

	"github.com/sonda2208/guardrails-challenge/model"
)

const (
	SecretKeyPattern = `(?i)\s*(\bBEGIN\b).*((PRIVATE KEY)|(PUBLIC KEY)\b)\s*`
)

type SecretKeyRule struct {
	re *regexp.Regexp
}

func NewSecretKeyRule() (*SecretKeyRule, error) {
	re, err := regexp.Compile(SecretKeyPattern)
	if err != nil {
		return nil, err
	}

	return &SecretKeyRule{
		re: re,
	}, nil
}

func (r SecretKeyRule) Name() string {
	return "SecretKey"
}

func (r SecretKeyRule) Metadata() model.FindingMetadata {
	return model.FindingMetadata{
		Description: "Use of hard-coded secret keys",
		Severity:    "HIGH",
	}
}

func (r SecretKeyRule) Scan(filePath string, fileBuf []byte) ([]*model.Finding, error) {
	positions := r.regexSearch(fileBuf)
	findings := make([]*model.Finding, len(positions))
	for i, p := range positions {
		findings[i] = &model.Finding{
			Type:     r.Name(),
			Metadata: r.Metadata(),
			Location: model.FindingLocation{
				Path:      filePath,
				Positions: p,
			},
		}
	}

	return findings, nil
}

func (r SecretKeyRule) regexSearch(fileBuf []byte) []model.FindingPosition {
	var res []model.FindingPosition

	locs := r.re.FindAllIndex(fileBuf, -1)
	lastStart := 0
	lastLineNumber := 0
	lastMatchIndex := 0
	lastLineStartIndex := 0
	for _, match := range locs {
		start := match[0]
		lineStart := lastLineStartIndex
		if idx := bytes.LastIndex(fileBuf[lastStart:start], []byte{'\n'}); idx >= 0 {
			lineStart = lastStart + idx + 1
		}

		lineNumber, matchIndex := r.computeLineNumber(fileBuf, lastLineNumber, lastMatchIndex, lineStart, start)

		lastMatchIndex = matchIndex
		lastLineNumber = lineNumber

		res = append(res, model.FindingPosition{
			Begin: model.PositionIndex{
				Line: lineNumber + 1,
			},
		})
	}

	return res
}

func (r SecretKeyRule) computeLineNumber(fileBuf []byte, lastLineNumber, lastMatchIndex, lineStart, matchStart int) (lineNumber, matchIndex int) {
	lineNumber = lastLineNumber + bytes.Count(fileBuf[lastMatchIndex:matchStart], []byte{'\n'})
	return lineNumber, lineStart
}
