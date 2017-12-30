package eep

import (
	"fmt"
	"strconv"
)

type scanner struct {
	source   string
	runes    []rune
	tokens   []*Token
	start    int
	current  int
	line     int
	keywords map[string]tokenKind
	errors   []error
}

func newScanner(s string) *scanner {
	return &scanner{
		source: s,
		runes:  []rune(s),
		tokens: []*Token{},
		line:   1,
		keywords: map[string]tokenKind{
			"and":   And,
			"false": False,
			"nil":   Nil,
			"or":    OR,
			"true":  True,
		},
	}
}

func (s *scanner) scanTokens() ([]*Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()

		if len(s.errors) > 0 {
			return s.tokens, s.errorsToError()
		}
	}
	s.tokens = append(
		s.tokens,
		&Token{
			TokenType: EOF,
			Lexeme:    "",
			Literal:   nil,
			Line:      s.line,
		},
	)
	return s.tokens, nil
}

func (s *scanner) scanToken() {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(LeftParent, nil)
	case ')':
		s.addToken(RightParent, nil)
	case '[':
		s.addToken(LeftBrace, nil)
	case ']':
		s.addToken(RightBrace, nil)
	case ',':
		s.addToken(Comma, nil)
	case '.':
		s.addToken(Dot, nil)
	case '-':
		s.addToken(Minus, nil)
	case '+':
		s.addToken(Plus, nil)
	case ';':
		s.addToken(Semicolon, nil)
	case '*':
		s.addToken(Star, nil)
	case '!':
		s.addToken(conditionalExp(s.match('='), BangEqual, Bang), nil)
	case '=':
		s.addToken(conditionalExp(s.match('='), EqualEqual, Equal), nil)
	case '<':
		s.addToken(conditionalExp(s.match('='), LessEqual, Less), nil)
	case '>':
		s.addToken(conditionalExp(s.match('='), GreaterEqual, Greater), nil)
	case '/':
		s.addToken(Slash, nil)
	case ' ', '\r', '\t': // 自动 break
	case '\n':
		s.line++
	case '"':
		s.getStr()
	default:
		if isDigits(c) {
			s.getNumber()
		} else if isAlpha(c) {
			s.getIdentifier()
		} else {
			s.errors = append(s.errors, newTokenError(s.line, "Unexpected Token!"))
		}
	}
}

func (s *scanner) advance() rune {
	s.current++
	return s.runes[s.current-1]
}

func (s *scanner) peek() rune {
	if s.current >= len(s.runes) {
		return '\000' // https://stackoverflow.com/questions/38007361/is-there-anyway-to-create-null-terminated-string-in-go
	}
	return s.runes[s.current]
}

func (s *scanner) peekNext() rune {
	if (s.current + 1) >= len(s.runes) {
		return '\000'
	}
	return s.runes[s.current+1]
}

func (s *scanner) addToken(tokenType tokenKind, literal interface{}) {
	text := string(s.runes[s.start:s.current])
	s.tokens = append(s.tokens, &Token{TokenType: tokenType, Lexeme: text, Literal: literal, Line: s.line})
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.runes)
}

func (s *scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.runes[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *scanner) getStr() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		s.errors = append(s.errors, newTokenError(s.line, "Unterminated string"))
		return
	}

	s.advance()
	value := string(s.runes[s.start+1 : s.current-1])
	s.addToken(String, value)
}

func (s *scanner) getNumber() {
	for isDigits(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigits(s.peekNext()) {
		s.advance()

		for isDigits(s.peek()) {
			s.advance()
		}
	}

	text := string(s.runes[s.start:s.current])
	number, _ := strconv.ParseFloat(text, 64)
	s.addToken(Number, number)
}

func (s *scanner) getIdentifier() {
	for isAlphaNumberic(s.peek()) {
		s.advance()
	}

	text := string(s.runes[s.start:s.current])
	tokenType, ok := s.keywords[text]
	if !ok {
		tokenType = Identifier
	}
	s.addToken(tokenType, nil)
}

func (s *scanner) errorsToError() error {
	var msg string
	if len(s.errors) > 0 {
		for _, e := range s.errors {
			msg += e.Error()
		}

		return fmt.Errorf(msg)
	}

	return nil
}
