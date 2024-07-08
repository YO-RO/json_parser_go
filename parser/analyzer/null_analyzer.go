package analyzer

import "regexp"

func mayBeNull(str string, i int) bool {
	return str[i] == 'n'
}

func extractNullAsToken(str string, startIdx int) (Token, int, error) {
	re := regexp.MustCompile(`^null\b`)
	loc := re.FindStringIndex(str[startIdx:])
	if loc == nil {
		return Token{}, 0, ErrUndefinedSymbol
	}
	endIdx := startIdx + loc[1]
	token := Token{
		Type:  Null,
		Value: nil,
	}
	return token, endIdx, nil
}
