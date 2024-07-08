package analyzer

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

func mayBeNumber(str string, i int) bool {
	matched, _ := regexp.MatchString(`-|\d`, str[i:i+1])
	return matched
}

func extractNumberAsToken(str string, startIdx int) (ValueToken, int, error) {
	re := regexp.MustCompile(`^-?(0|[1-9]\d*)(\.\d+)?(e[+-]?\d+)?`)
	loc := re.FindStringSubmatchIndex(str[startIdx:])
	if loc == nil {
		return ValueToken{}, 0, ErrUndefinedSymbol
	}
	endIdx := startIdx + loc[1]
	numStr := str[startIdx:endIdx]

	// 指数表現はstrconv.ParseFloat()でしか変換できない
	value, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		m := fmt.Sprintf("value: strconv.ParseFloat(%q, 64): ", numStr) +
			err.Error()
		panic(m)
	}

	var token ValueToken
	if value == math.Round(value) {
		token = ValueToken{
			tokenType: Int,
			value:     int(value),
		}
	} else {
		token = ValueToken{
			tokenType: Float,
			value:     value,
		}
	}
	return token, endIdx, nil
}
