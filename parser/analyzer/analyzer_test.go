package analyzer_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/YO-RO/mini-parser-go/parser/analyzer"
)

type testCase struct {
	name      string
	input     string
	expectErr error
	want      []analyzer.Token
}

func runTestCases(t *testing.T, tests []testCase) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := analyzer.Analyze([]byte(tc.input))
			if !errors.Is(tc.expectErr, err) {
				t.Errorf("err(%v) expects to be %v", err, tc.expectErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("analyzer.Analyze([]byte(%q)) == %v, want %v", tc.input, got, tc.want)
			}
		})
	}

}

func intToken(t *testing.T, value int) analyzer.Token {
	t.Helper()
	return analyzer.Token{
		Type:  analyzer.Int,
		Value: value,
	}
}

func floatToken(t *testing.T, value float64) analyzer.Token {
	t.Helper()
	return analyzer.Token{
		Type:  analyzer.Float,
		Value: value,
	}
}

func strToken(t *testing.T, value string) analyzer.Token {
	t.Helper()
	return analyzer.Token{
		Type:  analyzer.String,
		Value: value,
	}
}

func boolToken(t *testing.T, value bool) analyzer.Token {
	t.Helper()
	return analyzer.Token{
		Type:  analyzer.Bool,
		Value: value,
	}
}

func nullToken(t *testing.T) analyzer.Token {
	t.Helper()
	return analyzer.Token{
		Type:  analyzer.Null,
		Value: nil,
	}
}

func markToken(t *testing.T, value string) analyzer.Token {
	t.Helper()
	switch value {
	case `[`:
		return analyzer.Token{Type: analyzer.LeftSquareBracket, Value: `[`}
	case `]`:
		return analyzer.Token{Type: analyzer.RightSquareBracket, Value: `]`}
	case `{`:
		return analyzer.Token{Type: analyzer.LeftCurlyBracket, Value: `{`}
	case `}`:
		return analyzer.Token{Type: analyzer.RightCurlyBracket, Value: `}`}
	case `:`:
		return analyzer.Token{Type: analyzer.Colon, Value: `:`}
	case `,`:
		return analyzer.Token{Type: analyzer.Comma, Value: `,`}
	default:
		return analyzer.Token{}
	}
}

func TestNullAnalyzer(t *testing.T) {
	tests := []testCase{
		{
			"null",
			"null",
			nil,
			[]analyzer.Token{nullToken(t)},
		},
		{
			"nullでない文字[nulll]",
			"nulll",
			analyzer.ErrUndefinedSymbol,
			nil,
		},
		{
			"nullでない文字[nnull]",
			"nnull",
			analyzer.ErrUndefinedSymbol,
			nil,
		},
	}
	runTestCases(t, tests)
}

func TestBoolAnalyzer(t *testing.T) {
	tests := []testCase{
		{
			"真",
			"true",
			nil,
			[]analyzer.Token{boolToken(t, true)},
		},
		{
			"偽",
			"false",
			nil,
			[]analyzer.Token{boolToken(t, false)},
		},
		{
			"真(連続して記号が来る場合)",
			"true,",
			nil,
			[]analyzer.Token{
				boolToken(t, true),
				markToken(t, `,`),
			},
		},
		{
			"定義されていないキーワード[trueeeee]",
			"trueeeee",
			analyzer.ErrUndefinedSymbol,
			nil,
		},
		{
			"定義されていないキーワード[atrue]",
			"atrue",
			analyzer.ErrUndefinedSymbol,
			nil,
		},
		{
			"定義されていないキーワード[afalse]",
			"afalse",
			analyzer.ErrUndefinedSymbol,
			nil,
		},
		{
			"定義されていないキーワード[true111]",
			"true111",
			analyzer.ErrUndefinedSymbol,
			nil,
		},
	}
	runTestCases(t, tests)
}

func TestIntAnalyzer(t *testing.T) {
	tests := []testCase{
		{
			"ゼロ",
			"0",
			nil,
			[]analyzer.Token{intToken(t, 0)},
		},
		{
			"一つの整数",
			"123",
			nil,
			[]analyzer.Token{intToken(t, 123)},
		},
		{
			"複数の整数",
			"123 456 789",
			nil,
			[]analyzer.Token{
				intToken(t, 123),
				intToken(t, 456),
				intToken(t, 789),
			},
		},
		{
			"負の値",
			"-15",
			nil,
			[]analyzer.Token{intToken(t, -15)},
		},
		{
			"指数表現(符号なし)",
			"1e10",
			nil,
			[]analyzer.Token{intToken(t, 1e10)},
		},
		{
			"指数表現(符号あり)",
			"1e+10",
			nil,
			[]analyzer.Token{intToken(t, 1e+10)},
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
			[]analyzer.Token{floatToken(t, 12.34)},
		},
		{
			"複数のfloat",
			"12.34 56.78 90.12",
			nil,
			[]analyzer.Token{
				floatToken(t, 12.34),
				floatToken(t, 56.78),
				floatToken(t, 90.12),
			},
		},
		{
			"ゼロ点",
			"0.1",
			nil,
			[]analyzer.Token{floatToken(t, 0.1)},
		},
		{
			"-ゼロ点",
			"-0.1",
			nil,
			[]analyzer.Token{floatToken(t, -0.1)},
		},
		{
			"負の値",
			"-15.6",
			nil,
			[]analyzer.Token{floatToken(t, -15.6)},
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
			[]analyzer.Token{
				strToken(t, ""),
			},
		},
		{
			"文字列",
			`"string"`,
			nil,
			[]analyzer.Token{
				strToken(t, "string"),
			},
		},
		{
			"エスケープ文字付き文字列",
			`"\n, \\, \", \\n, \\\", \\"`,
			nil,
			[]analyzer.Token{
				strToken(t, "\n, \\, \", \\n, \\\", \\"),
			},
		},
		{
			"複数の文字列",
			`"hello""world" "foo"`,
			nil,
			[]analyzer.Token{
				strToken(t, "hello"),
				strToken(t, "world"),
				strToken(t, "foo"),
			},
		},
		{
			"不正な制御文字",
			`"line1
			line2"`,
			analyzer.ErrSyntax,
			nil,
		},
		{
			"文字列を閉じるための\"がない",
			`"string`,
			analyzer.ErrSyntax,
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
			[]analyzer.Token{
				markToken(t, `,`),
			},
		},
		{
			"コロン",
			`:`,
			nil,
			[]analyzer.Token{
				markToken(t, `:`),
			},
		},
		{
			"左中カッコ (left square bracket)",
			`[`,
			nil,
			[]analyzer.Token{
				markToken(t, `[`),
			},
		},
		{
			"右中カッコ (right square bracket)",
			`]`,
			nil,
			[]analyzer.Token{
				markToken(t, `]`),
			},
		},
		{
			"左大カッコ (left curly bracket)",
			`{`,
			nil,
			[]analyzer.Token{
				markToken(t, `{`),
			},
		},
		{
			"右大カッコ (right curly bracket)",
			`}`,
			nil,
			[]analyzer.Token{
				markToken(t, `}`),
			},
		},
		{
			"複数の記号",
			`{[:,]}`,
			nil,
			[]analyzer.Token{
				markToken(t, `{`),
				markToken(t, `[`),
				markToken(t, `:`),
				markToken(t, `,`),
				markToken(t, `]`),
				markToken(t, `}`),
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
			[]analyzer.Token{
				analyzer.Token(markToken(t, `{`)),

				strToken(t, "title"),
				analyzer.Token(markToken(t, `:`)),
				strToken(t, "go"),
				analyzer.Token(markToken(t, `,`)),

				strToken(t, "published"),
				analyzer.Token(markToken(t, `:`)),
				boolToken(t, true),
				analyzer.Token(markToken(t, `,`)),

				strToken(t, "year"),
				analyzer.Token(markToken(t, `:`)),
				intToken(t, 2025),
				analyzer.Token(markToken(t, `,`)),

				strToken(t, "rate"),
				analyzer.Token(markToken(t, `:`)),
				floatToken(t, 0.1),
				analyzer.Token(markToken(t, `,`)),

				strToken(t, "authors"),
				analyzer.Token(markToken(t, `:`)),
				analyzer.Token(markToken(t, `[`)),
				strToken(t, "ab"),
				analyzer.Token(markToken(t, `,`)),
				strToken(t, "a=b"),
				analyzer.Token(markToken(t, `,`)),
				strToken(t, "\"quotation gg\""),
				analyzer.Token(markToken(t, `]`)),
				analyzer.Token(markToken(t, `,`)),

				strToken(t, "desc"),
				analyzer.Token(markToken(t, `:`)),
				strToken(t, "This book is written about go language.\nBy gophers."),

				analyzer.Token(markToken(t, `}`)),
			},
		},
	}
	runTestCases(t, tests)
}
