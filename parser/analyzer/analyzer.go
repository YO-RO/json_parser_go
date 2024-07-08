package analyzer

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
	Null

	LeftSquareBracket
	RightSquareBracket
	LeftCurlyBracket
	RightCurlyBracket
	Colon
	Comma
)

type Token struct {
	Type  TokenType
	Value any
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

func Analyze(d []byte) ([]Token, error) {
	res := []Token{}

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
		case mayBeNumber(inputStr, i):
			token, endIdx, err := extractNumberAsToken(inputStr, i)
			if err != nil {
				return nil, err
			}
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
		case mayBeNull(inputStr, i):
			token, endIdx, err := extractNullAsToken(inputStr, i)
			if err != nil {
				return nil, err
			}
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
