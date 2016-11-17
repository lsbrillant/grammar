package grammar

import "testing"

var grammarValidityTests = [...]struct {
	grammar Grammar
	isvalid bool
}{
	{
		Grammar{
			map[Nonterminal][]Rule{
				Nonterminal("S"): {{"S", []Symbol{Terminal("a"), Nonterminal("B")}}},
				Nonterminal("B"): {{"B", []Symbol{Terminal("b")}}},
			},
		},
		true,
	},
	{
		Grammar{
			map[Nonterminal][]Rule{
				Nonterminal("S"): {{"S", []Symbol{Terminal("a"), Nonterminal("B")}}},
			},
		},
		false,
	},
	{
		Grammar{
			map[Nonterminal][]Rule{
				Nonterminal("A"): {{"A", []Symbol{Terminal("a"), Nonterminal("B")}}},
				Nonterminal("B"): {{"B", []Symbol{Terminal("b")}}},
			},
		},
		false,
	},
}

func TestIsValid(t *testing.T) {
	for i, test := range grammarValidityTests {
		if test.grammar.IsValid() != test.isvalid {
			t.Error("Error in test %d : grammer does not match expected", i)
		}
	}
}
