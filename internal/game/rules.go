package game

type Rules struct {
	StartingCards int
	HookRule      bool
	TrumpOrder    []Suit
}

func DefaultRules() Rules {
	return Rules{
		StartingCards: 7,
		HookRule:      true,
		TrumpOrder:    []Suit{Spades, Diamonds, Clubs, Hearts},
	}
}
