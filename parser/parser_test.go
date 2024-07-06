package parser_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/YO-RO/mini-parser-go/parser"
)

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

func TestIntAnalyzer(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr error
		want      []parser.Tokener
	}{
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

func TestAnalyzer(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr error
		want      []parser.Tokener
	}{
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
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parser.Analyze([]byte(tc.input))
			if !errors.Is(tc.expectErr, err) {
				t.Errorf("err(%s) expects to be %s", err, tc.expectErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("parser.Analyze([]byte(%s)) == %s, want %s", tc.input, got, tc.want)
			}
		})
	}
}
