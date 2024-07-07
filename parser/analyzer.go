package parser

import (
	"errors"
	"regexp"
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
