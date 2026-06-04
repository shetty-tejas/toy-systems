package main

import (
	"strings"
	"unicode"
)

type lexemekind int

type lexeme struct {
	kind  lexemekind
	value string
}

type stateFunc func(*lexer) stateFunc

const (
	paren lexemekind = iota
	name
	number
	stringl
	invalid
)

type lexer struct {
	source     []rune
	checkpoint int
	cursor     int

	lexemes []lexeme
}

func NewLexer(source string) *lexer {
	return &lexer{
		source:  []rune(source),
		lexemes: make([]lexeme, 0, 512),
	}
}

func (l *lexer) Lex() {
	for state := defaultState(l); state != nil; {
		state = state(l)
	}
}

func (l *lexer) emit(kind lexemekind) {
	if l.lexemeFound() {
		l.lexemes = append(l.lexemes, lexeme{kind: kind, value: string(l.source[l.checkpoint:l.cursor])})
	}

	l.checkpoint = l.cursor
}

func (l *lexer) lexemeFound() bool {
	return l.checkpoint != l.cursor
}

// State Functions

func defaultState(l *lexer) stateFunc {
	if l.cursor >= len(l.source) {
		return nil
	}

	ch := l.source[l.cursor]

	if isParen(ch) {
		l.cursor++
		l.emit(paren)

		return defaultState
	} else if unicode.IsSpace(ch) {
		return spaceState
	} else if unicode.IsDigit(ch) {
		return numberState
	} else if unicode.IsLetter(ch) || isOperator(ch) {
		return nameState
	} else if ch == '"' {
		return stringState
	}

	return invalidState
}

func spaceState(l *lexer) stateFunc {
	if l.cursor >= len(l.source) {
		return nil
	}

	if ch := l.source[l.cursor]; unicode.IsSpace(ch) {
		l.cursor++
		l.checkpoint = l.cursor
	}

	return defaultState
}

func nameState(l *lexer) stateFunc {
	if l.cursor >= len(l.source) {
		l.emit(name)
		return nil
	}

	if ch := l.source[l.cursor]; unicode.IsLetter(ch) || isOperator(ch) || unicode.IsDigit(ch) {
		l.cursor++

		return nameState
	} else {
		l.emit(name)

		return defaultState
	}
}

func numberState(l *lexer) stateFunc {
	if l.cursor >= len(l.source) {
		l.emit(number)
		return nil
	}

	if ch := l.source[l.cursor]; unicode.IsDigit(ch) {
		l.cursor++

		return numberState
	} else {
		l.emit(number)

		if unicode.IsSpace(ch) || isParen(ch) {
			return defaultState
		}

		return invalidState
	}
}

func invalidState(l *lexer) stateFunc {
	if l.cursor >= len(l.source) {
		l.emit(invalid)

		return nil
	}

	if ch := l.source[l.cursor]; unicode.IsSpace(ch) || isParen(ch) {
		l.emit(invalid)

		return defaultState
	}

	l.cursor++

	return invalidState
}

func stringState(l *lexer) stateFunc {
	if l.cursor >= len(l.source) {
		if l.cursor > len(l.source) {
			l.cursor = len(l.source)
		}

		return invalidState
	}

	if ch := l.source[l.cursor]; (l.lexemeFound() && ch != '"') || (!l.lexemeFound() && ch == '"') {
		if ch == '\\' {
			l.cursor++
		}

		l.cursor++

		return stringState
	} else {
		l.cursor++
		l.emit(stringl)

		return defaultState
	}
}

func isOperator(ch rune) bool {
	return strings.ContainsRune("+-*/=<>!%^&|", ch)
}

func isParen(ch rune) bool {
	return ch == '(' || ch == ')'
}
