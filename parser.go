package eep

import (
	"fmt"
)

type parser struct {
	tokens  []*Token
	current int
}

type parserErr struct {
	msg string
}

func (p parserErr) Error() string {
	return p.msg
}

func newParser(tokens []*Token) *parser {
	return &parser{tokens: tokens}
}

func (p *parser) parse() (exp Expr, err error) {
	defer func() {
		if e := recover(); e != nil {
			switch err1 := e.(type) {
			case string:
				err = fmt.Errorf(err1)
			case parserErr:
				err = err1
			}
		}
	}()

	return p.expression(), nil
}

func (p *parser) expression() Expr {
	return p.or()
}

func (p *parser) or() Expr {
	exp := p.and()

	for p.match(OR) {
		op := p.previous()
		right := p.and()

		exp = NewLogical(exp, right, op)
	}

	return exp
}

func (p *parser) and() Expr {
	exp := p.equality()

	for p.match(And) {
		op := p.previous()
		right := p.equality()

		exp = NewLogical(exp, right, op)
	}

	return exp
}

func (p *parser) equality() Expr {
	exp := p.comparison()

	for p.match(BangEqual, EqualEqual) {
		op := p.previous()
		right := p.comparison()
		exp = NewBinary(exp, right, op)
	}
	return exp
}

func (p *parser) comparison() Expr {
	exp := p.term()

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		op := p.previous()
		right := p.term()
		exp = NewBinary(exp, right, op)
	}
	return exp
}

func (p *parser) term() Expr {
	exp := p.factor()

	for p.match(Minus, Plus) {
		op := p.previous()
		right := p.factor()
		exp = NewBinary(exp, right, op)
	}

	return exp
}

func (p *parser) factor() Expr {
	exp := p.unary()

	for p.match(Slash, Star) {
		op := p.previous()
		right := p.unary()
		exp = NewBinary(exp, right, op)
	}

	return exp
}

func (p *parser) unary() Expr {
	if p.match(Bang, Minus) {
		op := p.previous()
		right := p.unary()
		return NewUnary(op, right)
	}

	return p.call()
}

func (p *parser) call() Expr {
	exp := p.primary()

	for {
		if p.match(LeftParent) {
			exp = p.finishCall(exp)
		} else {
			break
		}
	}

	return exp
}

func (p *parser) finishCall(callee Expr) Expr {
	var arguments []Expr

	if !p.check(RightParent) {
		if len(arguments) > 8 {
			panic(&parserErr{msg: "can not have more than 8 arguments"})
		}

		arguments = append(arguments, p.expression())

		for p.match(Comma) {
			arguments = append(arguments, p.expression())
		}
	}

	paren := p.consume(RightParent, "Expect `)` after arguments!")
	return NewCall(callee, paren, arguments)
}

func (p *parser) primary() Expr {
	if p.match(False) {
		return NewLiteral(false)
	}

	if p.match(True) {
		return NewLiteral(true)
	}

	if p.match(Nil) {
		return NewLiteral(nil)
	}

	if p.match(Number, String) {
		return NewLiteral(p.previous().Literal)
	}

	if p.match(Identifier) {
		return NewVariable(p.previous())
	}

	if p.match(LeftParent) {
		exp := p.expression()
		p.consume(RightParent, "Expect ')' after expression")
		return NewGrouping(exp)
	}
	panic(fmt.Sprintf("expect expression: %s", p.peek().Lexeme))
}

func (p *parser) consume(t tokenKind, msg string) *Token {
	if p.check(t) {
		return p.advance()
	}
	panic(&parserErr{msg: msg})
}

func (p *parser) match(tokenTypes ...tokenKind) bool {
	for _, t := range tokenTypes {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) check(t tokenKind) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == t
}

func (p *parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *parser) previous() *Token {
	return p.tokens[p.current-1]
}
