package analyzer

import (
	"regexp"
	"strconv"
)

func mayBeString(str string, i int) bool {
	return str[i] == '"'
}

func extractStringAsToken(str string, startIdx int) (ValueToken, int, error) {
	re := regexp.MustCompile(`^"[^\\]*?(\\.[^\\]*?)*?"`)
	loc := re.FindStringIndex(str[startIdx:])
	if loc == nil {
		return ValueToken{}, 0, ErrSyntax
	}
	endIdx := startIdx + loc[1]

	value, err := strconv.Unquote(str[startIdx:endIdx])
	if err != nil {
		return ValueToken{}, 0, ErrSyntax
	}
	token := ValueToken{
		tokenType: String,
		value:     value,
	}
	return token, endIdx, nil
}
