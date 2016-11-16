package grammar

import (
	"fmt"
	"os"
)

type Parser struct {
	scanner Scanner

	pos Position
	tok token
	lit string

	err        ErrorHandler
	ErrorCount int
}

func (p *Parser) Init(src []byte, e ErrorHandler) {
	p.scanner = NewScanner(src)
	p.err = e

	p.ErrorCount = 0

	p.pos, p.tok, p.lit = p.scanner.Scan()
}
func (p *Parser) error(msg string) {
	p.ErrorCount++
	if p.ErrorCount > 5 {
		fmt.Fprintln(os.Stdout, "to many errors")
		os.Exit(1)
	}
	p.err(p.pos, msg)
}

func (p *Parser) match(t token) (lit string) {
	if p.tok != t {
		p.error(fmt.Sprintf("expecting %s found %s", t, p.tok))
	}
	lit = p.lit
	p.next()
	return
}
func (p *Parser) next() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
}
func (p *Parser) passWhiteSpace() {
	for p.tok == Space {
		p.match(Space)
	}
}

func ParseGrammar(src []byte) (g Grammar) {
	p := Parser{}
	p.Init(src, func(pos Position, msg string) {
		// TODO make error handler not repeat errors
		fmt.Fprintf(os.Stderr, "error at %s %s\n", pos, msg)
	})
	for p.tok != Eof {
		r := Rule{}
		// Gobble up any irelivent whitespace
		p.passWhiteSpace()
		from := Nonterminal(p.match(Identifier))
		r.From = from
		p.passWhiteSpace()
		p.match(Arrow)
	rule:
		p.passWhiteSpace()
		for p.tok != Eof && p.tok != Space {
			switch p.tok {
			case Identifier:
				r.append(Nonterminal(p.match(Identifier)))
			case Literal:
				r.append(Terminal(p.match(Literal)))
			default:
				p.error(fmt.Sprintf("expecting Identifier or Literal found %s", p.tok))
			}
		}
		p.passWhiteSpace()
		g.AddRule(r)
		// goto implemeting a sort of do while looking thing
		// TODO dont use goto
		if p.tok == Pipe {
			p.match(Pipe)
			r = Rule{}
			r.From = from
			goto rule
		}
	}
	return g
}
