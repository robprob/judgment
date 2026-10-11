package game

// Associate each "play" in the stack with the player who placed it
type Play struct {
	PlayerIndex int
	Card        Card
}

type Trick struct {
	Plays []Play
	Trump Suit
}

// Declare the current winning PlayerID of the Trick
func (t Trick) Winner() int {
	if len(t.Plays) == 0 {
		return -1
	}
	winningPlay := t.Plays[0]
	lead := winningPlay.Card.Suit

	for _, challengePlay := range t.Plays[1:] {
		if beats(challengePlay.Card, winningPlay.Card, lead, t.Trump) {
			winningPlay = challengePlay
		}
	}
	return winningPlay.PlayerIndex
}

// Helper function determining whether passed card outright beats a comparison card
func beats(card, comp Card, lead, trump Suit) bool {
	if card.Suit != lead && card.Suit != trump {
		// Card is not lead or trump suit, cannot ever win
		return false
	}
	if comp.Suit != lead && comp.Suit != trump {
		// Comparison card is not lead or trump suit, passed card always wins
		return true
	}
	if card.Suit == comp.Suit {
		// Cards are equal suit and can be compared fairly based on rank
		return card.Rank > comp.Rank
	}
	// One is lead suit, other is trump suit (winner)
	return card.Suit == trump
}
