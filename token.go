package eep

import "fmt"

type tokenKind int

const (
	LeftParent tokenKind = iota  // (
	RightParent                  // )

	LeftBrace                    // [
	RightBrace                   // ]

	Comma                        // ,
	Dot                          // .
	Plus                         // +
	Minus  // -
	Star   // *
	Slash  // /
	Semicolon  // ;

	Bang  // !
	BangEqual  // !=
	Equal  //  =
	EqualEqual  // ==
	Greater     // >
	GreaterEqual   // >=
	Less           // <
	LessEqual      // <=

	Identifier
	String
	Number

	// KeyWords
	And
	False
	Nil
	OR
	True

	EOF
)

// Token lexeme
type Token struct {
	TokenType tokenKind
	Lexeme    string
	Literal   interface{}
	Line      int
}

// ToString token的字符串表示
func (tk *Token) ToString() string {
	return fmt.Sprintf("type: %d lexeme: %s literal: %v", int(tk.TokenType), tk.Lexeme, tk.Literal)
}