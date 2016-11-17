package grammar

type Grammar struct {
	Rules map[Nonterminal][]Rule
}

func (g *Grammar) AddRule(r Rule) {
	if g.Rules == nil {
		g.Rules = make(map[Nonterminal][]Rule)

	}
	if _, exists := g.Rules[r.From]; !exists {
		g.Rules[r.From] = make([]Rule, 0)
	}
	g.Rules[r.From] = append(g.Rules[r.From], r)
}

func (g *Grammar) IsValid() bool {
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

func (g *Grammar) IsRegular() bool {
	allLeft, allRight := true, true
	for _, rules := range g.Rules {
		for _, rule := range rules {
			if len(rule.Symbols) > 2 {
				return false
			}
			switch rule.Symbols[0].(type) {
			case Nonterminal:
				allRight = false
				if len(rule.Symbols) > 1 {
					switch rule.Symbols[1].(type) {
					case Nonterminal:
						return false
					}
				} else {
					return false
				}
			case Terminal:
				if len(rule.Symbols) > 1 {
					allLeft = false
					switch rule.Symbols[1].(type) {
					case Terminal:
						return false
					}
				}
			}
		}
	}
	return allLeft || allRight
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
