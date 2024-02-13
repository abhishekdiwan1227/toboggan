package ast

type Parser struct {
	tokens []SyntaxNode
	pos    int
}

func InitParser(code string) *Parser {
	lexer := InitLexer(code)

	var tokens []SyntaxNode
	for token := lexer.Lex(); token.Kind() != eof; token = lexer.Lex() {
		if token.Kind() == invalid {
			panic("Invalid token found")
		} else if token.Kind() != whitespace && token.Kind() != eof {
			tokens = append(tokens, token)
		}
	}

	return &Parser{tokens: tokens, pos: 0}
}

func (p Parser) peek(offset int) SyntaxNode {
	var index = p.pos + offset
	if index >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[index]
}

func (p Parser) current() SyntaxNode   { return p.peek(0) }
func (p Parser) lookahead() SyntaxNode { return p.peek(1) }

func (p *Parser) next() SyntaxNode {
	curr := p.current()
	p.pos++
	return curr
}

func (p *Parser) match(kind TokenKind) SyntaxNode {
	if p.current().Kind() == kind {
		return p.next()
	} else {
		panic("Invalid token found")
	}
}

func (p *Parser) Parse() SyntaxNode {
	return p.parseExpression(0)
}

func (p *Parser) parseExpression(parentPrecedence int) SyntaxNode {
	if p.current().Kind() == identifier && p.lookahead().Kind() == equals {
		identifier := p.primaryExpression()
		equals := p.match(equals)
		expression := p.parseExpression(0)
		return AssignmentExpresion{Expression: &Expression{SyntaxToken: &SyntaxToken{tokenKind: assignex}}, lhs: identifier, equalsToken: equals, rhs: expression}
	}

	var left SyntaxNode

	precedenceUnary := GetUnaryPrecedence(p.current().Kind())
	if precedenceUnary != 0 && precedenceUnary >= parentPrecedence {
		operator := p.next()
		operand := p.parseExpression(precedenceUnary)
		left = UnaryExpression{Expression: &Expression{SyntaxToken: &SyntaxToken{tokenKind: unaryex}}, operator: operator, operand: operand}
	} else {
		left = p.primaryExpression()
	}

	for {
		precedenceBinary := GetBinaryPrecedence(p.current().Kind())
		if precedenceBinary == 0 && precedenceBinary <= parentPrecedence {
			break
		}
		operator := p.next()
		right := p.parseExpression(precedenceBinary)
		left = BinaryExpression{Expression: &Expression{SyntaxToken: &SyntaxToken{tokenKind: binaryex}}, operator: operator, right: right, left: left}
	}
	return left
}

func (p *Parser) primaryExpression() SyntaxNode {
	if p.current().Kind() == openparan {
		left := p.next()
		expression := p.parseExpression(0)
		right := p.match(closeparan)
		return ParanthesesExpression{Expression: &Expression{SyntaxToken: &SyntaxToken{tokenKind: expression.Kind()}}, open: left, expression: expression, close: right}
	} else if p.current().Kind() == boolfalse || p.current().Kind() == booltrue {
		literal := p.next()
		return LiteralExpression{Expression: &Expression{SyntaxToken: &SyntaxToken{tokenKind: literal.Kind()}}, Value: literal.Kind() == booltrue, token: literal}
	} else if p.current().Kind() == identifier {
		name := p.next()
		return IdentifierExpression{Expression: &Expression{&SyntaxToken{tokenKind: nameex}}, name: name}
	}
	number := p.match(number).(SyntaxToken)
	return LiteralExpression{Expression: &Expression{SyntaxToken: &SyntaxToken{tokenKind: literal}}, Value: number.Value, token: number}
}
