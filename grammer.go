package grammer

type Grammer struct {
	Rules map[Nonterminal][]Rule
}

func (g *Grammer) AddRule(r Rule) {
	if g.Rules == nil {
		g.Rules = make(map[Nonterminal][]Rule)

	}
	if _, exists := g.Rules[r.From]; !exists {
		g.Rules[r.From] = make([]Rule, 0)
	}
	g.Rules[r.From] = append(g.Rules[r.From], r)
}

func (g *Grammer) IsValid() bool {
	if _, exists := g.Rules[StartSymbol]; !exists {
		return false
	}
	for _, rules := range g.Rules {
		for _, rule := range rules {
			for _, sym := range rule.Symbols {
				switch s := sym.(type) {
				case Nonterminal:
					if _, exists := g.Rules[s]; !exists {
						return false
					}
				case Terminal:
					// do nothing
				}
			}
		}
	}
	return true
}

type Rule struct {
	From    Nonterminal
	Symbols []Symbol
}

func (r *Rule) append(s Symbol) {
	if r.Symbols == nil {
		r.Symbols = make([]Symbol, 0)
	}
	r.Symbols = append(r.Symbols, s)
}

type Symbol interface {
	Sym()
}
type Terminal string
type Nonterminal string

var StartSymbol = Nonterminal("S")
var Lambda = Terminal("(lambda)")

func (s Terminal) Sym()    {}
func (s Nonterminal) Sym() {}
