package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type TokenType int

const (
	String TokenType = iota + 1
	Int
	Float
	Bool

	LeftSquareBracket
	RightSquareBracket
	LeftCurlyBracket
	RightCurlyBracket
	Colon
	Comma
)

var (
	ErrNoValue error = errors.New("no value")
)

// TokenはValueTokenとMarkToken
type Tokener interface {
	TokenType() TokenType
	Value() (any, error)
}

type ValueToken struct {
	tokenType TokenType
	value     any
}

func NewValueToken(tokenType TokenType, value any) ValueToken {
	return ValueToken{
		tokenType: tokenType,
		value:     value,
	}
}

func (vt ValueToken) TokenType() TokenType {
	return vt.tokenType
}

func (vt ValueToken) Value() (any, error) {
	return vt.value, nil
}

type MarkToken struct {
	tokenType TokenType
}

func NewMarkToken(tokenType TokenType) MarkToken {
	return MarkToken{tokenType: tokenType}
}

func (mt MarkToken) TokenType() TokenType {
	return mt.tokenType
}

func (mt MarkToken) Value() (any, error) {
	return nil, ErrNoValue
}

var (
	ErrSyntax          error = errors.New("invalid syntax")
	ErrUndefinedSymbol error = errors.New("undefined symbol")
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

func isNumber(str string, i int) bool {
	matched, _ := regexp.MatchString(`\d`, str[i:i+1])
	return matched
}

func mustExtractNumberAsToken(str string, startIdx int) (ValueToken, int) {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
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

func isSpace(str string, i int) bool {
	matched, _ := regexp.MatchString(`\s`, str[i:i+1])
	return matched
}

func skipSpaces(str string, startIdx int) int {
	re := regexp.MustCompile(`\s+`)
	loc := re.FindStringIndex(str[startIdx:])
	if loc == nil {
		return startIdx
	}
	return startIdx + loc[1]
}

func Analyze(d []byte) ([]Tokener, error) {
	res := []Tokener{}

	inputStr := string(d)
	for i := 0; i < len(inputStr); i++ {
		switch {
		case mayBeString(inputStr, i):
			token, endIdx, err := extractStringAsToken(inputStr, i)
			if err != nil {
				return nil, err
			}
			res = append(res, token)
			i = endIdx - 1
		case isNumber(inputStr, i):
			token, endIdx := mustExtractNumberAsToken(inputStr, i)
			res = append(res, token)
			i = endIdx - 1
		case mayBeBool(inputStr, i):
			token, endIdx, err := extractBoolAsToken(inputStr, i)
			if err != nil {
				return nil, err
			}
			res = append(res, token)
			i = endIdx - 1
		case isMark(inputStr, i):
			token, endIdx := mustExtractMark(inputStr, i)
			res = append(res, token)
			i = endIdx - 1
		case isSpace(inputStr, i):
			endIdx := skipSpaces(inputStr, i)
			i = endIdx - 1
		default:
			return nil, ErrUndefinedSymbol
		}
	}

	return res, nil
}
