package parser

import (
	"errors"
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

func Analyze(d []byte) ([]Tokener, error) {
	res := []Tokener{}

	inputStr := string(d)
	for i := 0; i < len(inputStr); i++ {
		switch inputStr[i] {
		case '"':
			// 直前に\がない"
			// または 直前に\が偶数回連続している"
			// `"`や`\\"`などがマッチ
			//
			// Capture Group 0 を全体のマッチだとすると、
			// Quotation Mark(")はCapture Group 5
			re := regexp.MustCompile(`(^|[^\\]|((^|[^\\])(\\\\)+))(")`)
			lastQuotationIdx :=
				re.FindStringSubmatchIndex(inputStr[i+1:])
			if len(lastQuotationIdx) == 0 {
				return nil, ErrSyntax
			}
			// idxsはinputStr[i+1]からのインデックスであるためi+1を足す
			beginIdx := i
			endIdx := i + 1 + lastQuotationIdx[11]

			value, err := strconv.Unquote(inputStr[beginIdx:endIdx])
			if err != nil {
				return nil, ErrSyntax
			}
			token := ValueToken{
				tokenType: String,
				value:     value,
			}
			res = append(res, token)

			i = endIdx - 1
		}
	}

	return res, nil
}
