package main

import (
	"bufio"
	"bytes"
	"io"
)

// Token is a token
type Token int

// Known tokens
const (
	ILLEGAL Token = iota
	EOF
	WHITESPACE
	NEWLINE
	NUMBER
	WORD
	PUNCT

	eof = rune(0)
)

// Tokens maps Token to name
var Tokens = map[Token]string{
	ILLEGAL:    "<illegal>",
	EOF:        "<eof>",
	WHITESPACE: "<whitespace>",
	NEWLINE:    "<newline>",
	NUMBER:     "<number>",
	WORD:       "<word>",
	PUNCT:      "<punctuation>",
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' }
func isNewline(ch rune) bool    { return ch == '\n' || ch == '\r' }
func isAlpha(ch rune) bool      { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }
func isNumeric(ch rune) bool    { return (ch >= '0' && ch <= '9') }
func isPunct(ch rune) bool {
	return (ch >= '!' && ch <= '/') ||
		(ch >= ':' && ch <= '@') ||
		(ch >= '[' && ch <= '`') ||
		(ch >= '{' && ch <= '~')
}
func isWordPunct(ch rune) bool {
	switch ch {
	case '_', '-', '+', '\'':
		return true
	default:
		return false
	}
}

// Scanner is a scanner
type Scanner struct {
	r *bufio.Reader
}

// NewScanner creates a scanner
func NewScanner(r io.Reader) *Scanner { return &Scanner{r: bufio.NewReader(r)} }

func (s *Scanner) read() rune {
	if ch, _, err := s.r.ReadRune(); err == nil {
		return ch
	}
	return eof
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return WHITESPACE, buf.String()
}

func (s *Scanner) scanNewlines() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isNewline(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return NEWLINE, buf.String()
}

func (s *Scanner) scanWord() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isNumeric(ch) && !isAlpha(ch) && !isWordPunct(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return WORD, buf.String()
}

func (s *Scanner) scanNumeric() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isNumeric(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return NUMBER, buf.String()
}

// Scan scans
func (s *Scanner) Scan() (tok Token, lit string) {
	ch := s.read()
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isNewline(ch) {
		s.unread()
		return s.scanNewlines()
	} else if isAlpha(ch) {
		s.unread()
		return s.scanWord()
	} else if isNumeric(ch) {
		s.unread()
		return s.scanNumeric()
	} else if isPunct(ch) {
		return PUNCT, string(ch)
	}

	switch ch {
	case eof:
		return EOF, ""
	}
	return ILLEGAL, string(ch)
}
