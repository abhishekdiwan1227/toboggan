//go:generate stringer -type=TokenKind

package ast

type TokenKind int

const (
	literal TokenKind = iota
	whitespace
	eof
	invalid
	plus
	dash
	star
	slash
	not
	number
	identifier
	binaryex
	unaryex
	paranthesesex
	booltrue
	boolfalse
	dequals
	notequals
	damp
	dpipe
	equals
	nameex
	assignex
	openparan
	closeparan
)

func GetUnaryPrecedence(kind TokenKind) int {
	if kind == dash || kind == plus {
		return 6
	} else {
		return 0
	}
}

func GetBinaryPrecedence(kind TokenKind) int {
	if kind == star || kind == slash {
		return 5
	} else if kind == plus || kind == dash {
		return 4
	} else if kind == dequals || kind == notequals {
		return 3
	} else if kind == damp {
		return 2
	} else if kind == dpipe {
		return 1
	} else {
		return 0
	}
}

func GetKeywordKind(raw string) TokenKind {
	if raw == "true" {
		return booltrue
	} else if raw == "false" {
		return boolfalse
	} else {
		return identifier
	}
}

type SyntaxNode interface {
	Kind() TokenKind
	Children() []SyntaxNode
	TokenValue() any
}

type SyntaxToken struct {
	tokenKind TokenKind
	pos       int
	text      string
	Value     any
}

func (t SyntaxToken) Kind() TokenKind        { return t.tokenKind }
func (t SyntaxToken) Children() []SyntaxNode { return nil }
func (t SyntaxToken) TokenValue() any        { return t.Value }

type Literal struct {
	*SyntaxToken
	Value any
}

func (l Literal) Children() []SyntaxNode { return []SyntaxNode{l.SyntaxToken} }

type Expression struct {
	*SyntaxToken
}

type BinaryExpression struct {
	*Expression
	operator SyntaxNode
	left     SyntaxNode
	right    SyntaxNode
}

func (e BinaryExpression) Children() []SyntaxNode { return []SyntaxNode{e.left, e.operator, e.right} }

type UnaryExpression struct {
	*Expression
	operator SyntaxNode
	operand  SyntaxNode
}

func (e UnaryExpression) Children() []SyntaxNode { return []SyntaxNode{e.operator, e.operand} }

type IdentifierExpression struct {
	*Expression
	name SyntaxNode
}

func (e IdentifierExpression) Children() []SyntaxNode { return []SyntaxNode{e.name} }

type ParanthesesExpression struct {
	*Expression
	open       SyntaxNode
	expression SyntaxNode
	close      SyntaxNode
}

func (e ParanthesesExpression) Children() []SyntaxNode {
	return []SyntaxNode{e.open, e.expression, e.close}
}

type LiteralExpression struct {
	*Expression
	Value any
	token SyntaxNode
}

func (e LiteralExpression) Children() []SyntaxNode { return []SyntaxNode{e.token} }

type AssignmentExpresion struct {
	*Expression
	lhs         SyntaxNode
	equalsToken SyntaxNode
	rhs         SyntaxNode
}

func (e AssignmentExpresion) Children() []SyntaxNode {
	return []SyntaxNode{e.lhs, e.equalsToken, e.rhs}
}
