package game

import "slices"

type Player struct {
	Name   string
	Hand   []Card
	Bid    int
	Tricks int
	Score  int
}

// LegalCards returns the valid playable cards in the player's hand for this trick
func (p Player) LegalCards(lead *Suit) []Card {
	if lead == nil {
		// No lead suit has been established, player chooses among any cards in hand
		return slices.Clone(p.Hand)
	}

	// Determine which, if any, cards in hand match the lead suit
	var matchLead []Card
	for _, card := range p.Hand {
		if card.Suit == *lead {
			matchLead = append(matchLead, card)
		}
	}
	if len(matchLead) > 0 {
		// Player must play one of their lead suit-matching cards
		return matchLead
	}
	// All cards valid
	return slices.Clone(p.Hand)
}
