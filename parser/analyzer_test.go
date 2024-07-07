package parser_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/YO-RO/mini-parser-go/parser"
)

type testCase struct {
	name      string
	input     string
	expectErr error
	want      []parser.Tokener
}

func runTestCases(t *testing.T, tests []testCase) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parser.Analyze([]byte(tc.input))
			if !errors.Is(tc.expectErr, err) {
				t.Errorf("err(%v) expects to be %v", err, tc.expectErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("parser.Analyze([]byte(%q)) == %v, want %v", tc.input, got, tc.want)
			}
		})
	}

}

func intToken(t *testing.T, value int) parser.ValueToken {
	t.Helper()
	return parser.NewValueToken(parser.Int, value)
}

func floatToken(t *testing.T, value float64) parser.ValueToken {
	t.Helper()
	return parser.NewValueToken(parser.Float, value)
}

func strToken(t *testing.T, value string) parser.ValueToken {
	t.Helper()
	return parser.NewValueToken(parser.String, value)
}

func boolToken(t *testing.T, value bool) parser.ValueToken {
	t.Helper()
	return parser.NewValueToken(parser.Bool, value)
}

func TestBoolAnalyzer(t *testing.T) {
	tests := []testCase{
		{
			"真",
			"true",
			nil,
			[]parser.Tokener{boolToken(t, true)},
		},
		{
			"偽",
			"false",
			nil,
			[]parser.Tokener{boolToken(t, false)},
		},
		{
			"真(連続して記号が来る場合)",
			"true,",
			nil,
			[]parser.Tokener{
				boolToken(t, true),
				parser.NewMarkToken(parser.Comma),
			},
		},
		{
			"定義されていないキーワード[trueeeee]",
			"trueeeee",
			parser.ErrUndefinedSymbol,
			nil,
		},
		{
			"定義されていないキーワード[atrue]",
			"atrue",
			parser.ErrUndefinedSymbol,
			nil,
		},
		{
			"定義されていないキーワード[afalse]",
			"afalse",
			parser.ErrUndefinedSymbol,
			nil,
		},
		{
			"定義されていないキーワード[true111]",
			"true111",
			parser.ErrUndefinedSymbol,
			nil,
		},
	}
	runTestCases(t, tests)
}

func TestIntAnalyzer(t *testing.T) {
	tests := []testCase{
		{
			"一つの整数",
			"123",
			nil,
			[]parser.Tokener{intToken(t, 123)},
		},
		{
			"複数の整数",
			"123 456 789",
			nil,
			[]parser.Tokener{
				intToken(t, 123),
				intToken(t, 456),
				intToken(t, 789),
			},
		},
	}
	runTestCases(t, tests)
}

func TestFloatAnalyzer(t *testing.T) {
	tests := []testCase{
		{
			"一つのfloat",
			"12.34",
			nil,
			[]parser.Tokener{floatToken(t, 12.34)},
		},
		{
			"複数のfloat",
			"12.34 56.78 90.12",
			nil,
			[]parser.Tokener{
				floatToken(t, 12.34),
				floatToken(t, 56.78),
				floatToken(t, 90.12),
			},
		},
	}
	runTestCases(t, tests)
}

func TestStringAnalyzer(t *testing.T) {
	tests := []testCase{
		{
			"空文字列",
			`""`,
			nil,
			[]parser.Tokener{
				strToken(t, ""),
			},
		},
		{
			"文字列",
			`"string"`,
			nil,
			[]parser.Tokener{
				strToken(t, "string"),
			},
		},
		{
			"エスケープ文字付き文字列",
			`"\n, \\, \", \\n, \\\", \\"`,
			nil,
			[]parser.Tokener{
				strToken(t, "\n, \\, \", \\n, \\\", \\"),
			},
		},
		{
			"複数の文字列",
			`"hello""world" "foo"`,
			nil,
			[]parser.Tokener{
				strToken(t, "hello"),
				strToken(t, "world"),
				strToken(t, "foo"),
			},
		},
		{
			"不正な制御文字",
			`"line1
			line2"`,
			parser.ErrSyntax,
			nil,
		},
		{
			"文字列を閉じるための\"がない",
			`"string`,
			parser.ErrSyntax,
			nil,
		},
	}
	runTestCases(t, tests)
}

func TestMarkAnalyzer(t *testing.T) {
	tests := []testCase{
		{
			"コンマ",
			`,`,
			nil,
			[]parser.Tokener{
				parser.NewMarkToken(parser.Comma),
			},
		},
		{
			"コロン",
			`:`,
			nil,
			[]parser.Tokener{
				parser.NewMarkToken(parser.Colon),
			},
		},
		{
			"左中カッコ (left square bracket)",
			`[`,
			nil,
			[]parser.Tokener{
				parser.NewMarkToken(parser.LeftSquareBracket),
			},
		},
		{
			"右中カッコ (right square bracket)",
			`]`,
			nil,
			[]parser.Tokener{
				parser.NewMarkToken(parser.RightSquareBracket),
			},
		},
		{
			"左大カッコ (left curly bracket)",
			`{`,
			nil,
			[]parser.Tokener{
				parser.NewMarkToken(parser.LeftCurlyBracket),
			},
		},
		{
			"右大カッコ (right curly bracket)",
			`}`,
			nil,
			[]parser.Tokener{
				parser.NewMarkToken(parser.RightCurlyBracket),
			},
		},
		{
			"複数の記号",
			`{[:,]}`,
			nil,
			[]parser.Tokener{
				parser.NewMarkToken(parser.LeftCurlyBracket),
				parser.NewMarkToken(parser.LeftSquareBracket),
				parser.NewMarkToken(parser.Colon),
				parser.NewMarkToken(parser.Comma),
				parser.NewMarkToken(parser.RightSquareBracket),
				parser.NewMarkToken(parser.RightCurlyBracket),
			},
		},
	}
	runTestCases(t, tests)
}

func TestJson(t *testing.T) {
	tests := []testCase{
		{
			"jsonの例",
			`
			{
				"title": "go",
				"published": true,
				"year": 2025,
				"rate": 0.1,
				"authors": [ "ab", "a=b", "\"quotation gg\"" ],
				"desc": "This book is written about go language.\nBy gophers."
			}
			`,
			nil,
			[]parser.Tokener{
				parser.NewMarkToken(parser.LeftCurlyBracket),

				strToken(t, "title"),
				parser.NewMarkToken(parser.Colon),
				strToken(t, "go"),
				parser.NewMarkToken(parser.Comma),

				strToken(t, "published"),
				parser.NewMarkToken(parser.Colon),
				boolToken(t, true),
				parser.NewMarkToken(parser.Comma),

				strToken(t, "year"),
				parser.NewMarkToken(parser.Colon),
				intToken(t, 2025),
				parser.NewMarkToken(parser.Comma),

				strToken(t, "rate"),
				parser.NewMarkToken(parser.Colon),
				floatToken(t, 0.1),
				parser.NewMarkToken(parser.Comma),

				strToken(t, "authors"),
				parser.NewMarkToken(parser.Colon),
				parser.NewMarkToken(parser.LeftSquareBracket),
				strToken(t, "ab"),
				parser.NewMarkToken(parser.Comma),
				strToken(t, "a=b"),
				parser.NewMarkToken(parser.Comma),
				strToken(t, "\"quotation gg\""),
				parser.NewMarkToken(parser.RightSquareBracket),
				parser.NewMarkToken(parser.Comma),

				strToken(t, "desc"),
				parser.NewMarkToken(parser.Colon),
				strToken(t, "This book is written about go language.\nBy gophers."),

				parser.NewMarkToken(parser.RightCurlyBracket),
			},
		},
	}
	runTestCases(t, tests)
}
