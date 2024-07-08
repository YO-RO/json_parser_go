package analyzer

import "regexp"

func extractBoolAsToken(str string, startIdx int) (Token, int, error) {
	re := regexp.MustCompile(`^(true|false)\b`)
	loc := re.FindStringIndex(str[startIdx:])
	if loc == nil {
		return Token{}, 0, errNoMatch
	}
	endIdx := startIdx + loc[1]

	var value bool
	if str[startIdx:endIdx] == "true" {
		value = true
	} else {
		value = false
	}

	token := Token{
		Type:  Bool,
		Value: value,
	}
	return token, endIdx, nil
}
