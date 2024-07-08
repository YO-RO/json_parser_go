package analyzer

import (
	"regexp"
)

func extractMarkAsToken(str string, startIdx int) (Token, int, error) {
	re := regexp.MustCompile(`^[,:\[\]{}]`)
	mark := re.FindString(str[startIdx:])
	if mark == "" {
		return Token{}, 0, errNoMatch
	}
	endIdx := startIdx + 1 // markは1文字
	switch mark {
	case ",":
		return Token{Type: Comma, Value: `,`}, endIdx, nil
	case ":":
		return Token{Type: Colon, Value: `:`}, endIdx, nil
	case "[":
		return Token{Type: LeftSquareBracket, Value: `[`}, endIdx, nil
	case "]":
		return Token{Type: RightSquareBracket, Value: `]`}, endIdx, nil
	case "{":
		return Token{Type: LeftCurlyBracket, Value: `{`}, endIdx, nil
	case "}":
		return Token{Type: RightCurlyBracket, Value: `}`}, endIdx, nil
	default:
		panic("out of range: mark must match one of ,:[]{} : mark: " + mark)
	}
}
