package main

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []lexeme
	}{
		{
			name:     "empty input",
			input:    "",
			expected: []lexeme{},
		},
		{
			name:     "only whitespace",
			input:    " \t\n\r ",
			expected: []lexeme{},
		},
		{
			name:  "parens",
			input: "()",
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "simple expression",
			input: "(add 2 3)",
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "add"},
				{kind: number, value: "2"},
				{kind: number, value: "3"},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "nested expression",
			input: "(add 2 (subtract 4 2))",
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "add"},
				{kind: number, value: "2"},
				{kind: paren, value: "("},
				{kind: name, value: "subtract"},
				{kind: number, value: "4"},
				{kind: number, value: "2"},
				{kind: paren, value: ")"},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "names with various letters",
			input: "foo BAR baz",
			expected: []lexeme{
				{kind: name, value: "foo"},
				{kind: name, value: "BAR"},
				{kind: name, value: "baz"},
			},
		},
		{
			name:  "signs as operator names (add)",
			input: "(+ 1 2)",
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "+"},
				{kind: number, value: "1"},
				{kind: number, value: "2"},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "signs as operator names (subtract)",
			input: "(- 4 2)",
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "-"},
				{kind: number, value: "4"},
				{kind: number, value: "2"},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "signs as operator names (nested)",
			input: "(+ (- 3 1) 2)",
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "+"},
				{kind: paren, value: "("},
				{kind: name, value: "-"},
				{kind: number, value: "3"},
				{kind: number, value: "1"},
				{kind: paren, value: ")"},
				{kind: number, value: "2"},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "invalid characters",
			input: "(add ! @)",
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "add"},
				{kind: name, value: "!"},
				{kind: invalid, value: "@"},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "mixed continuous without spaces",
			input: "(add2)",
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "add2"},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "string literal support (super tiny compiler feature)",
			input: `(concat "hello" "world")`,
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "concat"},
				{kind: stringl, value: `"hello"`},
				{kind: stringl, value: `"world"`},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "string literal with escaped quotes",
			input: `(print "hello \"world\"")`,
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "print"},
				{kind: stringl, value: `"hello \"world\""`},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "mixed continuous without spaces (numbers first)",
			input: "(2add)",
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: number, value: "2"},
				{kind: invalid, value: "add"},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "multiple escaped quotes",
			input: `(print "she said \"hi\" and \"bye\"")`,
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "print"},
				{kind: stringl, value: `"she said \"hi\" and \"bye\""`},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "escaped quote at start",
			input: `(print "\"hello")`,
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "print"},
				{kind: stringl, value: `"\"hello"`},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "escaped quote at end",
			input: `(print "hello\"")`,
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "print"},
				{kind: stringl, value: `"hello\""`},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "only escaped quotes",
			input: `(print "\"\"")`,
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "print"},
				{kind: stringl, value: `"\"\""`},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "escaped backslash before closing quote",
			input: `(print "path\\")`,
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "print"},
				{kind: stringl, value: `"path\\"`},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "empty string",
			input: `(print "")`,
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "print"},
				{kind: stringl, value: `""`},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "input ending with whitespace",
			input: "add ",
			expected: []lexeme{
				{kind: name, value: "add"},
			},
		},
		{
			name:  "input ending mid-name (EOF)",
			input: "foo",
			expected: []lexeme{
				{kind: name, value: "foo"},
			},
		},
		{
			name:  "input ending mid-number (EOF)",
			input: "123",
			expected: []lexeme{
				{kind: number, value: "123"},
			},
		},
		{
			name:  "single char name",
			input: "(a 1)",
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "a"},
				{kind: number, value: "1"},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "multiple invalid chars together",
			input: "(@@@)",
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: invalid, value: "@@@"},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "unicode letters in names",
			input: "(café 1)",
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "café"},
				{kind: number, value: "1"},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "deeply nested",
			input: "(a (b (c 1)))",
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "a"},
				{kind: paren, value: "("},
				{kind: name, value: "b"},
				{kind: paren, value: "("},
				{kind: name, value: "c"},
				{kind: number, value: "1"},
				{kind: paren, value: ")"},
				{kind: paren, value: ")"},
				{kind: paren, value: ")"},
			},
		},
		{
			name:  "unterminated string (EOF)",
			input: `(print "hello`,
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "print"},
				{kind: invalid, value: `"hello`},
			},
		},
		{
			name:  "dangling escape at EOF",
			input: `(print "hello\`,
			expected: []lexeme{
				{kind: paren, value: "("},
				{kind: name, value: "print"},
				{kind: invalid, value: `"hello\`},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			l.Lex()

			if len(tt.expected) == 0 && len(l.lexemes) == 0 {
				return
			}

			// Just to print output for failures clearly while debugging implementations
			// We'll use a slightly nicer formatting
			if !reflect.DeepEqual(l.lexemes, tt.expected) {
				t.Errorf("\nInput: %q\nExpected: %v\nGot:      %v", tt.input, tt.expected, l.lexemes)
			}
		})
	}
}
