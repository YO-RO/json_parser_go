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
	ErrSyntax error = errors.New("invalid syntax")
)

func analyzeStringToken(str string, firstQuotationIdx int) (ValueToken, int, error) {
	// 直前に\がない"
	// または 直前に\が偶数回連続している"
	// `"`や`\\"`などがマッチ
	//
	// Capture Group 0 を全体のマッチだとすると、
	// Quotation Mark(")はCapture Group 5
	re := regexp.MustCompile(`(^|[^\\]|((^|[^\\])(\\\\)+))(")`)
	matchedLastQuotationIdx :=
		re.FindStringSubmatchIndex(str[firstQuotationIdx+1:])
	if len(matchedLastQuotationIdx) == 0 {
		return ValueToken{}, 0, ErrSyntax
	}
	// idxsはstr[firstQuotationIdx+1]からのインデックスであるためfirstQuotationIdx+1を足す
	beginIdx := firstQuotationIdx
	endIdx := firstQuotationIdx + 1 + matchedLastQuotationIdx[11]

	value, err := strconv.Unquote(str[beginIdx:endIdx])
	if err != nil {
		return ValueToken{}, 0, ErrSyntax
	}
	token := ValueToken{
		tokenType: String,
		value:     value,
	}
	return token, endIdx, nil
}

func analyzeNumberToken(str string, firstNumberIdx int) (ValueToken, int, error) {
	var token ValueToken
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	loc := re.FindStringSubmatchIndex(str[firstNumberIdx:])
	// 小数点以降のマッチはloc[2]からloc[3]
	// loc[2] == -1 ならマッチしていないことになる
	if loc[2] == -1 {
		value, err := strconv.Atoi(str[firstNumberIdx : firstNumberIdx+loc[1]])
		if err != nil {
			// errorを埋め込みたくないため%v
			return ValueToken{}, 0, fmt.Errorf("internal error: %v", err)
		}
		token = ValueToken{
			tokenType: Int,
			value:     value,
		}
	} else {
		value, err := strconv.ParseFloat(str[firstNumberIdx:firstNumberIdx+loc[1]], 64)
		if err != nil {
			// errorを埋め込みたくないため%v
			return ValueToken{}, 0, fmt.Errorf("internal error: %v", err)
		}
		token = ValueToken{
			tokenType: Float,
			value:     value,
		}
	}
	return token, firstNumberIdx + loc[1], nil
}

func Analyze(d []byte) ([]Tokener, error) {
	res := []Tokener{}

	inputStr := string(d)
	for i := 0; i < len(inputStr); i++ {
		switch inputStr[i] {
		case '"':
			token, endIdx, err := analyzeStringToken(inputStr, i)
			if err != nil {
				return nil, err
			}
			res = append(res, token)
			i = endIdx - 1
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			token, endIdx, err := analyzeNumberToken(inputStr, i)
			if err != nil {
				return nil, err
			}
			res = append(res, token)
			i = endIdx - 1
		}
	}

	return res, nil
}
