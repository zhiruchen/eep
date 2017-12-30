package eep

import (
	"unicode"
)

// ConditionalExp 三元表达式
func conditionalExp(condition bool, v1, v2 tokenKind) tokenKind {
	if condition {
		return v1
	}
	return v2
}

func isDigits(r rune) bool {
	return unicode.IsDigit(r)
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isAlphaNumberic(c rune) bool {
	return isAlpha(c) || isDigits(c)
}
