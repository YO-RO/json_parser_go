package analyzer

import (
	"regexp"
	"strconv"
)

func mayBeString(str string, i int) bool {
	return str[i] == '"'
}

func extractStringAsToken(str string, startIdx int) (Token, int, error) {
	re := regexp.MustCompile(`^"[^\\]*?(\\.[^\\]*?)*?"`)
	loc := re.FindStringIndex(str[startIdx:])
	if loc == nil {
		return Token{}, 0, ErrSyntax
	}
	endIdx := startIdx + loc[1]

	value, err := strconv.Unquote(str[startIdx:endIdx])
	if err != nil {
		return Token{}, 0, ErrSyntax
	}
	token := Token{
		Type:  String,
		Value: value,
	}
	return token, endIdx, nil
}
