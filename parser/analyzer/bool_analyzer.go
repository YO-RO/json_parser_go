package analyzer

import "regexp"

func mayBeBool(str string, i int) bool {
	return str[i] == 't' || str[i] == 'f'
}

func extractBoolAsToken(str string, startIdx int) (ValueToken, int, error) {
	// ?: はグループをキャプチャしない
	re := regexp.MustCompile(`(true|false)(?:[\s,:"{}\[\]]|$)`)
	loc := re.FindStringSubmatchIndex(str[startIdx:])
	if loc == nil {
		return ValueToken{}, 0, ErrUndefinedSymbol
	}
	endIdx := startIdx + loc[3]

	var value bool
	if str[startIdx:endIdx] == "true" {
		value = true
	} else {
		value = false
	}

	token := ValueToken{
		tokenType: Bool,
		value:     value,
	}
	return token, endIdx, nil
}
