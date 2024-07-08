package analyzer

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

func extractNumberAsToken(str string, startIdx int) (Token, int, error) {
	re := regexp.MustCompile(`^-?(0|[1-9]\d*)(\.\d+)?(e[+-]?\d+)?`)
	loc := re.FindStringSubmatchIndex(str[startIdx:])
	if loc == nil {
		return Token{}, 0, errNoMatch
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

	var token Token
	if value == math.Round(value) {
		token = Token{
			Type:  Int,
			Value: int(value),
		}
	} else {
		token = Token{
			Type:  Float,
			Value: value,
		}
	}
	return token, endIdx, nil
}
