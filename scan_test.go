package grammar

import "testing"

var scanAndParseTests = [...]struct {
	program string
	tokens  []token
	grammar Grammar
}{
	{
		`S -> ab`,
		[]token{
			Identifier,
			Space,
			Arrow,
			Space,
			Literal,
			Literal,
		},
		Grammar{
			map[Nonterminal][]Rule{
				Nonterminal("S"): {{"S", []Symbol{Terminal("a"), Terminal("b")}}},
			},
		},
	},
	{
		`S -> aB
		 # comment
		 B -> b`,
		[]token{
			Identifier,
			Space,
			Arrow,
			Space,
			Literal,
			Identifier,
			Space,
			Identifier,
			Space,
			Arrow,
			Space,
			Literal,
		},
		Grammar{
			map[Nonterminal][]Rule{
				Nonterminal("S"): {{"S", []Symbol{Terminal("a"), Nonterminal("B")}}},
				Nonterminal("B"): {{"B", []Symbol{Terminal("b")}}},
			},
		},
	},
	{
		`S -> aB | Ba
		 B -> b`,
		[]token{
			Identifier,
			Space,
			Arrow,
			Space,
			Literal,
			Identifier,
			Space,
			Pipe,
			Space,
			Identifier,
			Literal,
			Space,
			Identifier,
			Space,
			Arrow,
			Space,
			Literal,
		},
		Grammar{
			map[Nonterminal][]Rule{
				Nonterminal("S"): {
					{"S", []Symbol{Terminal("a"), Nonterminal("B")}},
					{"S", []Symbol{Nonterminal("B"), Terminal("a")}},
				},
				Nonterminal("B"): {{"B", []Symbol{Terminal("b")}}},
			},
		},
	},
}

func TestScanner(t *testing.T) {
	for i, test := range scanAndParseTests {
		s := NewScanner([]byte(test.program))
		for j, expected := range test.tokens {
			pos, tok, lit := s.Scan()
			if tok != expected {
				t.Errorf("Error in test %d : token %d is %s should be %s at %s is %s",
					i, j, tok, expected, pos, lit)
			}
		}
	}
}

func TestParse(t *testing.T) {
	for i, test := range scanAndParseTests {
		g := ParseGrammar([]byte(test.program))
		for key, rules := range g.Rules {
			for k, rule := range rules {
				if rule.From != key {
					t.Errorf("Error in test %d : key %s does not match %s", i, rule.From, key)
				}
				if rule.From != test.grammar.Rules[key][k].From {
					t.Errorf("Error in test %d : rule.From %s does not match expected %s",
						i, rule.From, test.grammar.Rules[key][k].From)
				}
				for h, sym := range rule.Symbols {
					if sym != test.grammar.Rules[key][k].Symbols[h] {
						t.Errorf("Error in test %d : symbol %s does not match expected %s",
							i, sym, test.grammar.Rules[key][k].Symbols[h])
					}
				}
			}
		}
	}
}
