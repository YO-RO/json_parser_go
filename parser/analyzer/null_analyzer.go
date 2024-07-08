package analyzer

import "regexp"

func extractNullAsToken(str string, startIdx int) (Token, int, error) {
	re := regexp.MustCompile(`^null\b`)
	loc := re.FindStringIndex(str[startIdx:])
	if loc == nil {
		return Token{}, 0, errNoMatch
	}
	endIdx := startIdx + loc[1]
	token := Token{
		Type:  Null,
		Value: nil,
	}
	return token, endIdx, nil
}
