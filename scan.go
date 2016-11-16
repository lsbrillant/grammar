package grammer

import (
	"fmt"
	"os"
	"strconv"
	//	"unicode"
	"unicode/utf8"
)

// Position type taken from go/token
type Position struct {
	Offset int // offset, starting at 0
	Line   int // line number, starting at 1
	Column int // column number, starting at 1 (byte count)
}

func (pos *Position) IsValid() bool { return pos.Line > 0 }

// String returns a string in one of two forms:
//
//	line:column         valid position without file name
//	-                   invalid position without file name
//
func (pos Position) String() (s string) {
	if pos.IsValid() {
		s = fmt.Sprintf("%d:%d", pos.Line, pos.Column)
	} else {
		s = "-"
	}
	return
}

type token int

const (
	Identifier token = iota
	Literal
	Arrow
	Pipe

	Space
	Eof
)

var tokens = [...]string{
	Identifier: "IDENTIFIER",
	Literal:    "Literal",
	Arrow:      "->",
	Pipe:       "|",
	Space:      "SPACE",
	Eof:        "<<EOF>>",
}

func (tok token) String() string {
	s := ""
	if 0 <= tok && tok < token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

type ErrorHandler func(pos Position, msg string)
type Scanner struct {
	src []byte
	err ErrorHandler

	ch         rune
	offset     int
	rdOffset   int
	lineCount  int
	lineOffset int

	lit string

	ErrorCount int
}

const bom = 0xFEFF // byte order mark, only permitted as very first character

// from go/scanner with minor tweaks
func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		if s.ch == '\n' {
			s.lineCount++
			s.lineOffset = s.offset
		}
		r, w := rune(s.src[s.rdOffset]), 1
		switch {
		case r == 0:
			s.error("illegal character NUL")
		case r >= utf8.RuneSelf:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				s.error("illegal UTF-8 encoding")
			} else if r == bom && s.offset > 0 {
				s.error("illegal byte order mark")
			}
		}
		s.rdOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.lineCount++
			s.lineOffset = s.offset
		}
		s.ch = -1 // eof
	}
}
func (s *Scanner) Init(src []byte, err ErrorHandler) {
	s.src = src
	s.err = err

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.lineCount = 1
	s.lineOffset = 0
	s.ErrorCount = 0

	s.next()
	if s.ch == bom {
		s.next() // ignore BOM at file beginning
	}
}
func (s *Scanner) AtEof() bool {
	return s.ch == -1
}

func (s *Scanner) Pos() Position {
	return Position{s.offset, s.lineCount, (s.offset - s.lineOffset)}
}

func (s *Scanner) error(msg string) {
	s.ErrorCount++
	s.err(s.Pos(), msg)
}

func (s *Scanner) skipWhitespace() {
	for isWhiteSpace(s.ch) {
		s.next()
	}
}
func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
func isLowerCase(ch rune) bool {
	return 'a' <= ch && ch <= 'z'
}
func isUpperCase(ch rune) bool {
	return 'A' <= ch && ch <= 'Z'
}
func (s *Scanner) Scan() (pos Position, tok token, lit string) {
scanAgain:

	pos = s.Pos()
	switch ch := s.ch; {
	case isWhiteSpace(ch):
		tok = Space
		s.skipWhitespace()
	case isLowerCase(ch):
		tok = Literal
		lit = string(ch)
		s.next()
	case isUpperCase(ch):
		tok = Identifier
		offs := s.offset
		s.next()
		for ch = s.ch; '0' <= ch && ch <= '9'; ch = s.ch {
			s.next()
		}
		lit = string(s.src[offs:s.offset])
	default:
		s.next()
		switch ch {
		case -1:
			tok = Eof
		case '-':
			if s.ch != '>' {
				s.error(fmt.Sprintf("Expexting '->' found '-%c'", s.ch))
			}
			s.next()
			tok = Arrow
		case '|':
			tok = Pipe
		case '#':
			// comment
			for s.ch != '\n' && s.ch >= 0 {
				s.next()
			}
			s.skipWhitespace()
			goto scanAgain
		}
	}
	return
}

func NewScanner(src []byte) Scanner {
	var scanner Scanner
	var defaultErrorHandler ErrorHandler
	defaultErrorHandler = func(pos Position, msg string) {
		fmt.Fprintf(os.Stderr, "error at %s %s\n", pos, msg)
	}
	scanner.Init(src, defaultErrorHandler)
	return scanner
}
