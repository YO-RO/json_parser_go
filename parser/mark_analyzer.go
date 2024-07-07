package parser

import (
	"fmt"
	"regexp"
)

func isMark(str string, i int) bool {
	matched, _ := regexp.MatchString(`[,:\[\]{}]`, str[i:i+1])
	return matched
}

func mustExtractMark(str string, startIdx int) (MarkToken, int) {
	re := regexp.MustCompile(`[,:\[\]{}]`)
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
		return MarkToken{tokenType: Comma}, endIdx
	case ":":
		return MarkToken{tokenType: Colon}, endIdx
	case "[":
		return MarkToken{tokenType: LeftSquareBracket}, endIdx
	case "]":
		return MarkToken{tokenType: RightSquareBracket}, endIdx
	case "{":
		return MarkToken{tokenType: LeftCurlyBracket}, endIdx
	case "}":
		return MarkToken{tokenType: RightCurlyBracket}, endIdx
	default:
		panic("out of range: mark must match one of ,:[]{} : mark: " + mark)
	}
}
