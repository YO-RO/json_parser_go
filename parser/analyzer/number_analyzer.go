package analyzer

import (
	"fmt"
	"regexp"
	"strconv"
)

func isNumber(str string, i int) bool {
	matched, _ := regexp.MatchString(`\d`, str[i:i+1])
	return matched
}

func mustExtractNumberAsToken(str string, startIdx int) (ValueToken, int) {
	re := regexp.MustCompile(`^\d+(\.\d+)?`)
	loc := re.FindStringSubmatchIndex(str[startIdx:])
	if loc == nil {
		m := fmt.Sprintf(
			"loc must not be nil: loc: re.FindStringSubmatchIndex(%q)",
			str[startIdx:],
		)
		panic(m)
	}

	var token ValueToken
	endIdx := startIdx + loc[1]
	numStr := str[startIdx:endIdx]
	// 小数点以降のマッチはloc[2]からloc[3]
	// loc[2] == -1 ならマッチしていないことになる
	if loc[2] == -1 /*整数なら*/ {
		value, err := strconv.Atoi(numStr)
		if err != nil {
			m := fmt.Sprintf("value: strconv.Atoi(%q): ", numStr) + err.Error()
			panic(m)
		}
		token = ValueToken{
			tokenType: Int,
			value:     value,
		}
	} else /*floatなら*/ {
		value, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			m := fmt.Sprintf("value: strconv.ParseFloat(%q, 64): ", numStr) +
				err.Error()
			panic(m)
		}
		token = ValueToken{
			tokenType: Float,
			value:     value,
		}
	}
	return token, endIdx
}
