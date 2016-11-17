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
			t.Errorf("Error in test %d : grammer does not match expected", i)
		}
	}
}

var grammarHeirarchyTests = [...]struct {
	grammar   Grammar
	isRegular bool
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
				Nonterminal("S"): {{"S", []Symbol{Nonterminal("B"), Terminal("a")}}},
				Nonterminal("B"): {{"B", []Symbol{Terminal("b")}}},
			},
		},
		true,
	},
	{
		Grammar{
			map[Nonterminal][]Rule{
				Nonterminal("S"): {{"S", []Symbol{Nonterminal("B"), Terminal("a"), Terminal("c")}}},
				Nonterminal("B"): {{"B", []Symbol{Terminal("b")}}},
			},
		},
		false,
	},
	{
		Grammar{
			map[Nonterminal][]Rule{
				Nonterminal("S"): {{"S", []Symbol{Nonterminal("B"), Nonterminal("B")}}},
				Nonterminal("B"): {{"B", []Symbol{Terminal("b")}}},
			},
		},
		false,
	},
	{
		Grammar{
			map[Nonterminal][]Rule{
				Nonterminal("S"): {
					{"S", []Symbol{Terminal("a"), Nonterminal("B")}},
					{"S", []Symbol{Terminal("b"), Nonterminal("A")}},
				},
				Nonterminal("A"): {{"A", []Symbol{Terminal("a"), Nonterminal("B")}}},
				Nonterminal("B"): {{"B", []Symbol{Terminal("b")}}},
			},
		},
		true,
	},
	{
		Grammar{
			map[Nonterminal][]Rule{
				Nonterminal("S"): {
					{"S", []Symbol{Terminal("a"), Nonterminal("B")}},
					{"S", []Symbol{Terminal("b"), Nonterminal("A")}},
				},
				Nonterminal("A"): {{"A", []Symbol{Nonterminal("B"), Terminal("a")}}},
				Nonterminal("B"): {{"B", []Symbol{Terminal("b")}}},
			},
		},
		false,
	},
}

func TestIsRegular(t *testing.T) {
	for i, test := range grammarHeirarchyTests {
		if test.grammar.IsRegular() != test.isRegular {
			t.Errorf("Error in test %d : grammer does not match expected", i)
		}
	}
}
