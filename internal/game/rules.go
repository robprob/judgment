package game

type Rules struct {
	StartingCards     int
	HookRule          bool
	TrumpOrder        []Suit
	BlindOneCardRound bool
}

func DefaultRules() Rules {
	return Rules{
		StartingCards:     7,
		HookRule:          true,
		TrumpOrder:        []Suit{Spades, Diamonds, Clubs, Hearts},
		BlindOneCardRound: true,
	}
}
