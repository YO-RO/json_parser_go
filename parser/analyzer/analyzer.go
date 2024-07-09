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

type extractTokenFunc func(string, int) (Token, int, error)

var (
	ErrSyntax          error = errors.New("invalid syntax")
	ErrUndefinedSymbol error = errors.New("undefined symbol")

	errNoMatch error = errors.New("no match")
)

func skipSpaces(str string, startIdx int) (int, bool) {
	re := regexp.MustCompile(`^\s+`)
	loc := re.FindStringIndex(str[startIdx:])
	if loc == nil {
		return 0, false
	}
	return startIdx + loc[1], true
}

func Analyze(d []byte) ([]Token, error) {
	res := []Token{}
	inputStr := string(d)

	extractTokenFuncs := []extractTokenFunc{
		extractStringAsToken,
		extractNumberAsToken,
		extractBoolAsToken,
		extractMarkAsToken,
		extractNullAsToken,
	}
	for i := 0; i < len(inputStr); {
		// skipする時にcontinueしないと
		// inputStrが空白のみの時にpanic(out of range)になる
		endIdx, skip := skipSpaces(inputStr, i)
		if skip {
			i = endIdx
			continue
		}

		matched := false
		for _, extract := range extractTokenFuncs {
			token, endIdx, err := extract(inputStr, i)
			if err != nil {
				if errors.Is(errNoMatch, err) {
					continue
				}
				return nil, err
			}
			matched = true
			res = append(res, token)
			i = endIdx
			break
		}
		if !matched {
			return nil, ErrUndefinedSymbol
		}
	}

	return res, nil
}
