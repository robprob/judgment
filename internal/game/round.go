package game

import "fmt"

type Round struct {
	Players        []Player
	CardsPerPlayer int
	Trick          Trick
}

// Run a single round of Judgment
func (r *Round) Run() {
	// Instantiate a new deck
	deck := NewDeck()
	// Shuffle and deal cards to each player
	deck.Shuffle()
	for i := range r.Players {
		r.Players[i].Hand = deck.DrawN(r.CardsPerPlayer)
		// DEBUG: Print players/hands
		fmt.Printf("Player %s: %v\n", r.Players[i].Name, r.Players[i].Hand)
	}

	// Players play through each trick
	// TEMP: Randomly play one of valid hand cards

	// Calculate scores
}
