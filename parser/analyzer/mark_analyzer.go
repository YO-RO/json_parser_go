package analyzer

import (
	"fmt"
	"regexp"
)

func isMark(str string, i int) bool {
	matched, _ := regexp.MatchString(`[,:\[\]{}]`, str[i:i+1])
	return matched
}

func mustExtractMark(str string, startIdx int) (Token, int) {
	re := regexp.MustCompile(`^[,:\[\]{}]`)
	mark := re.FindString(str[startIdx:])
	if mark == "" {
		m := fmt.Sprintf(
			"mark must not be empty: mark: re.FindString(%q)",
			str[startIdx:],
		)
		panic(m)
	}
	endIdx := startIdx + 1 // markは1文字
	switch mark {
	case ",":
		return Token{Type: Comma, Value: `,`}, endIdx
	case ":":
		return Token{Type: Colon, Value: `:`}, endIdx
	case "[":
		return Token{Type: LeftSquareBracket, Value: `[`}, endIdx
	case "]":
		return Token{Type: RightSquareBracket, Value: `]`}, endIdx
	case "{":
		return Token{Type: LeftCurlyBracket, Value: `{`}, endIdx
	case "}":
		return Token{Type: RightCurlyBracket, Value: `}`}, endIdx
	default:
		panic("out of range: mark must match one of ,:[]{} : mark: " + mark)
	}
}
