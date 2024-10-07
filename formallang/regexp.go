package formallang

// RegExp - basic struct for regular expression
type RegExp struct {
	abc  map[rune]struct{}
	tree regExpNode
}

// Token - basic struct that string must be sliced
type Token struct {
	Symb          rune
	Servicable bool
}

// ToString - convert to 
func (reg RegExp) ToString() string {
	return reg.tree.ToString(lowPriority)
}

// RegExpFromTokens - construct regular expression from string
func RegExpFromTokens(tokens []Token) (*RegExp, error) {
	dict := make(map[rune]struct{})

	for _, token := range tokens {
		if !token.Servicable {
			dict[token.Symb] = struct{}{}
		}
	}

	return RegExpFromTokensWithDict(tokens, dict)
}

// RegExpFromTokensWithDict - construct regular expression from string with given alphabet
func RegExpFromTokensWithDict(tokens []Token, abc map[rune]struct{}) (*RegExp, error) {
	regexpnode, err := createRegExpNodes(tokens)

	res := &RegExp{
		abc: abc,
		tree: regexpnode,
	}

	return res, err
}

// Test - for test
func Test(tokens []Token) string {
	expr, _ := RegExpFromTokens(tokens)
	return expr.ToString()
}
