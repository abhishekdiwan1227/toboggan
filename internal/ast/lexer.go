package ast

import (
	"strconv"
	"unicode"
)

type Lexer struct {
	pos  int
	code string
}

func InitLexer(code string) *Lexer { return &Lexer{pos: 0, code: code} }

func (l Lexer) peek(offset int) rune {
	index := l.pos + offset
	if index >= len(l.code) {
		return '\000'
	}
	return rune(l.code[index])
}

func (l Lexer) current() rune   { return l.peek(0) }
func (l Lexer) lookahead() rune { return l.peek(1) }

func (l *Lexer) Lex() SyntaxNode {
	if l.current() == '\000' {
		return SyntaxToken{tokenKind: eof, pos: l.pos, Value: '\000'}
	}

	start := l.pos

	if l.current() == ' ' {
		for l.current() == ' ' {
			l.pos++
		}

		return SyntaxToken{tokenKind: whitespace, pos: start, Value: " "}
	}

	if l.current() == '+' {
		l.pos++
		return SyntaxToken{tokenKind: plus, pos: start, Value: "+"}
	} else if l.current() == '-' {
		l.pos++
		return SyntaxToken{tokenKind: dash, pos: start, Value: "-"}
	} else if l.current() == '*' {
		l.pos++
		return SyntaxToken{tokenKind: star, pos: start, Value: "*"}
	} else if l.current() == '/' {
		l.pos++
		return SyntaxToken{tokenKind: slash, pos: start, Value: "/"}
	} else if l.current() == '!' {
		if l.lookahead() == '=' {
			l.pos += 2
			return SyntaxToken{tokenKind: notequals, pos: start, Value: "!="}
		}
		l.pos++
		return SyntaxToken{tokenKind: not, pos: start, Value: "!"}
	} else if l.current() == '=' {
		if l.lookahead() == '=' {
			l.pos += 2
			return SyntaxToken{tokenKind: dequals, pos: start, Value: "=="}
		}
		l.pos++
		return SyntaxToken{tokenKind: equals, pos: start, Value: "="}
	} else if l.current() == '&' && l.lookahead() == '&' {
		l.pos += 2
		return SyntaxToken{tokenKind: damp, pos: start, Value: "&&"}
	} else if l.current() == '|' && l.lookahead() == '|' {
		l.pos += 2
		return SyntaxToken{tokenKind: dpipe, pos: start, Value: "||"}
	} else if l.current() == '(' {
		l.pos++
		return SyntaxToken{tokenKind: openparan, pos: start, Value: "("}
	} else if l.current() == ')' {
		l.pos++
		return SyntaxToken{tokenKind: closeparan, pos: start, Value: ")"}
	}

	if unicode.IsNumber(l.current()) {
		isNumber := true
		for unicode.IsNumber(l.current()) {
			l.pos++
		}
		for unicode.IsLetter(l.current()) {
			isNumber = false
			l.pos++
		}

		raw := l.code[start:l.pos]

		if !isNumber {
			return SyntaxToken{tokenKind: identifier, pos: start, text: raw, Value: raw}
		}

		value, err := strconv.Atoi(raw)
		if err != nil {
			panic(err)
		}
		return SyntaxToken{tokenKind: number, pos: start, text: raw, Value: value}
	}

	if unicode.IsLetter(l.current()) {
		for unicode.IsLetter(l.current()) || unicode.IsNumber(l.current()) {
			l.pos++
		}
		raw := l.code[start:l.pos]

		kind := GetKeywordKind(raw)

		return SyntaxToken{tokenKind: kind, pos: start, text: raw, Value: raw}
	}

	return SyntaxToken{tokenKind: invalid, pos: l.pos, text: string(l.current())}
}
