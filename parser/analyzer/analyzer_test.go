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
	want      []analyzer.Tokener
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

func intToken(t *testing.T, value int) analyzer.ValueToken {
	t.Helper()
	return analyzer.NewValueToken(analyzer.Int, value)
}

func floatToken(t *testing.T, value float64) analyzer.ValueToken {
	t.Helper()
	return analyzer.NewValueToken(analyzer.Float, value)
}

func strToken(t *testing.T, value string) analyzer.ValueToken {
	t.Helper()
	return analyzer.NewValueToken(analyzer.String, value)
}

func boolToken(t *testing.T, value bool) analyzer.ValueToken {
	t.Helper()
	return analyzer.NewValueToken(analyzer.Bool, value)
}

func nullToken(t *testing.T) analyzer.ValueToken {
	t.Helper()
	return analyzer.NewValueToken(analyzer.Null, nil)
}

func TestNullAnalyzer(t *testing.T) {
	tests := []testCase{
		{
			"null",
			"null",
			nil,
			[]analyzer.Tokener{nullToken(t)},
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
			[]analyzer.Tokener{boolToken(t, true)},
		},
		{
			"偽",
			"false",
			nil,
			[]analyzer.Tokener{boolToken(t, false)},
		},
		{
			"真(連続して記号が来る場合)",
			"true,",
			nil,
			[]analyzer.Tokener{
				boolToken(t, true),
				analyzer.NewMarkToken(analyzer.Comma),
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
			[]analyzer.Tokener{intToken(t, 0)},
		},
		{
			"一つの整数",
			"123",
			nil,
			[]analyzer.Tokener{intToken(t, 123)},
		},
		{
			"複数の整数",
			"123 456 789",
			nil,
			[]analyzer.Tokener{
				intToken(t, 123),
				intToken(t, 456),
				intToken(t, 789),
			},
		},
		{
			"負の値",
			"-15",
			nil,
			[]analyzer.Tokener{intToken(t, -15)},
		},
		{
			"指数表現(符号なし)",
			"1e10",
			nil,
			[]analyzer.Tokener{intToken(t, 1e10)},
		},
		{
			"指数表現(符号あり)",
			"1e+10",
			nil,
			[]analyzer.Tokener{intToken(t, 1e+10)},
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
			[]analyzer.Tokener{floatToken(t, 12.34)},
		},
		{
			"複数のfloat",
			"12.34 56.78 90.12",
			nil,
			[]analyzer.Tokener{
				floatToken(t, 12.34),
				floatToken(t, 56.78),
				floatToken(t, 90.12),
			},
		},
		{
			"ゼロ点",
			"0.1",
			nil,
			[]analyzer.Tokener{floatToken(t, 0.1)},
		},
		{
			"-ゼロ点",
			"-0.1",
			nil,
			[]analyzer.Tokener{floatToken(t, -0.1)},
		},
		{
			"負の値",
			"-15.6",
			nil,
			[]analyzer.Tokener{floatToken(t, -15.6)},
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
			[]analyzer.Tokener{
				strToken(t, ""),
			},
		},
		{
			"文字列",
			`"string"`,
			nil,
			[]analyzer.Tokener{
				strToken(t, "string"),
			},
		},
		{
			"エスケープ文字付き文字列",
			`"\n, \\, \", \\n, \\\", \\"`,
			nil,
			[]analyzer.Tokener{
				strToken(t, "\n, \\, \", \\n, \\\", \\"),
			},
		},
		{
			"複数の文字列",
			`"hello""world" "foo"`,
			nil,
			[]analyzer.Tokener{
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
			[]analyzer.Tokener{
				analyzer.NewMarkToken(analyzer.Comma),
			},
		},
		{
			"コロン",
			`:`,
			nil,
			[]analyzer.Tokener{
				analyzer.NewMarkToken(analyzer.Colon),
			},
		},
		{
			"左中カッコ (left square bracket)",
			`[`,
			nil,
			[]analyzer.Tokener{
				analyzer.NewMarkToken(analyzer.LeftSquareBracket),
			},
		},
		{
			"右中カッコ (right square bracket)",
			`]`,
			nil,
			[]analyzer.Tokener{
				analyzer.NewMarkToken(analyzer.RightSquareBracket),
			},
		},
		{
			"左大カッコ (left curly bracket)",
			`{`,
			nil,
			[]analyzer.Tokener{
				analyzer.NewMarkToken(analyzer.LeftCurlyBracket),
			},
		},
		{
			"右大カッコ (right curly bracket)",
			`}`,
			nil,
			[]analyzer.Tokener{
				analyzer.NewMarkToken(analyzer.RightCurlyBracket),
			},
		},
		{
			"複数の記号",
			`{[:,]}`,
			nil,
			[]analyzer.Tokener{
				analyzer.NewMarkToken(analyzer.LeftCurlyBracket),
				analyzer.NewMarkToken(analyzer.LeftSquareBracket),
				analyzer.NewMarkToken(analyzer.Colon),
				analyzer.NewMarkToken(analyzer.Comma),
				analyzer.NewMarkToken(analyzer.RightSquareBracket),
				analyzer.NewMarkToken(analyzer.RightCurlyBracket),
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
			[]analyzer.Tokener{
				analyzer.NewMarkToken(analyzer.LeftCurlyBracket),

				strToken(t, "title"),
				analyzer.NewMarkToken(analyzer.Colon),
				strToken(t, "go"),
				analyzer.NewMarkToken(analyzer.Comma),

				strToken(t, "published"),
				analyzer.NewMarkToken(analyzer.Colon),
				boolToken(t, true),
				analyzer.NewMarkToken(analyzer.Comma),

				strToken(t, "year"),
				analyzer.NewMarkToken(analyzer.Colon),
				intToken(t, 2025),
				analyzer.NewMarkToken(analyzer.Comma),

				strToken(t, "rate"),
				analyzer.NewMarkToken(analyzer.Colon),
				floatToken(t, 0.1),
				analyzer.NewMarkToken(analyzer.Comma),

				strToken(t, "authors"),
				analyzer.NewMarkToken(analyzer.Colon),
				analyzer.NewMarkToken(analyzer.LeftSquareBracket),
				strToken(t, "ab"),
				analyzer.NewMarkToken(analyzer.Comma),
				strToken(t, "a=b"),
				analyzer.NewMarkToken(analyzer.Comma),
				strToken(t, "\"quotation gg\""),
				analyzer.NewMarkToken(analyzer.RightSquareBracket),
				analyzer.NewMarkToken(analyzer.Comma),

				strToken(t, "desc"),
				analyzer.NewMarkToken(analyzer.Colon),
				strToken(t, "This book is written about go language.\nBy gophers."),

				analyzer.NewMarkToken(analyzer.RightCurlyBracket),
			},
		},
	}
	runTestCases(t, tests)
}
