package analyzer

import (
	"regexp"
	"strconv"
)

func mayBeString(str string, i int) bool {
	return str[i] == '"'
}

func extractStringAsToken(str string, startIdx int) (ValueToken, int, error) {
	// 直前に\がない"
	// または 直前に\が偶数回連続している"
	// `"`や`\\"`などがマッチ
	re := regexp.MustCompile(`(?:^|[^\\]|(?:(?:^|[^\\])(?:\\\\)+))(")`)
	loc := re.FindStringSubmatchIndex(str[startIdx+1:])
	if loc == nil {
		return ValueToken{}, 0, ErrSyntax
	}
	// idxsはstr[firstQuotationIdx+1]からのインデックスであるためfirstQuotationIdx+1を足す
	endIdx := startIdx + 1 + loc[3]

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
