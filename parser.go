package main

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type token struct {
	T Token
	V string
}

func (t token) String() string {
	return t.V
}

func (t token) DebugString() string {
	return fmt.Sprintf("%s: %s", Tokens[t.T], t.V)
}

type line []token

func (l line) String() string {
	var s []string
	for _, t := range l {
		s = append(s, t.String())
	}
	return strings.Join(s, " ")
}

func (l line) DebugString() string {
	var s []string
	for _, t := range l {
		s = append(s, t.DebugString())
	}
	return strings.Join(s, " ")
}

// Parser is a parser
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token
		lit string
		n   int
	}
}

// NewParser creates a Parser
func NewParser(r io.Reader) *Parser { return &Parser{s: NewScanner(r)} }

func (p *Parser) scan() (tok Token, lit string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	tok, lit = p.s.Scan()
	p.buf.tok, p.buf.lit = tok, lit
	return
}

func (p *Parser) unscan() { p.buf.n = 1 }

func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WHITESPACE {
		tok, lit = p.scan()
	}
	return
}

// Parse parses
func (p *Parser) Parse() error {
	var tokens line
	for {
		tok, lit := p.scanIgnoreWhitespace()
		switch tok {
		case EOF:
			handleLine(tokens)
			log.Printf("Done")
			return nil
		case ILLEGAL:
			log.Fatalf("Illegal Token: %v", lit)
		case NEWLINE:
			handleLine(tokens)
			tokens = nil
		default:
			tokens = append(tokens, token{T: tok, V: lit})
		}
	}
}

func genLine(tokens ...Token) line {
	var l line
	for _, tok := range tokens {
		l = append(l, token{T: tok})
	}
	return l
}

func genLineV(tokens ...token) line {
	var l line
	for _, tok := range tokens {
		l = append(l, tok)
	}
	return l
}

func matches(lhs, rhs line) bool {
	log.Printf("Comparing: \n%s\n%s", lhs.DebugString(), rhs.DebugString())
	if len(lhs) != len(rhs) {
		return false
	}
	for ix, l := range lhs {
		if l.T != rhs[ix].T {
			return false
		}
		if l.V != "" && l.V != rhs[ix].V {
			return false
		}
	}
	return true
}

var (
	timestamp = genLineV(
		token{T: PUNCT, V: "["},
		token{T: NUMBER},
		token{T: PUNCT, V: ":"},
		token{T: NUMBER},
		token{T: PUNCT, V: ":"},
		token{T: NUMBER},
		token{T: PUNCT, V: "]"})
	timerange = genLineV(
		token{T: NUMBER},
		token{T: PUNCT, V: "-"},
		token{T: NUMBER})
)

func handleLine(tokens line) {
	if len(tokens) >= 7 && matches(timestamp, tokens[:7]) {
		log.Printf("Timestamp Line")
	}
	log.Printf("%v", tokens)
}
