package analyzer

import "regexp"

func mayBeNull(str string, i int) bool {
	return str[i] == 'n'
}

func extractNullAsToken(str string, startIdx int) (ValueToken, int, error) {
	re := regexp.MustCompile(`^null\b`)
	loc := re.FindStringIndex(str[startIdx:])
	if loc == nil {
		return ValueToken{}, 0, ErrUndefinedSymbol
	}
	endIdx := startIdx + loc[1]
	token := ValueToken{
		tokenType: Null,
		value:     nil,
	}
	return token, endIdx, nil
}
