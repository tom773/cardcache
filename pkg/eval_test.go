package pkg

import (
	"testing"
)

func TestLex(t *testing.T) {
	lex([]byte("hello"))
}

func TestLex2(t *testing.T) {
	lex([]byte("[1, 2, 3]"))
}
