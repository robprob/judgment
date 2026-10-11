package game

import (
	"fmt"
	"math/rand/v2"
)

type Round struct {
	Players        []Player
	CardsPerPlayer int
	Trump          Suit
	Tricks         []Trick
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
		fmt.Printf("%s: %v\n", r.Players[i].Name, r.Players[i].Hand)
	}

	// TODO: Determine rotating leading player's position
	leaderIndex := 0
	leader := &r.Players[leaderIndex]
	// Players state bids
	// TEMP: Random int [0, min(3, CardsPerPlayer)] capped at CardsPerPlayer
	// TODO: Enforce bid "hook" rule (will be on/true by default)
	for i := range r.Players {
		r.Players[i].Bid = rand.IntN(min(3, r.CardsPerPlayer) + 1)
	}

	// Players play through each trick

	// DEBUG: Single trick
	trick := Trick{Trump: r.Trump}
	// Make leading play, determining suit and logging to current trick
	leadCard := leader.Hand[0]
	if !leader.Play(leadCard) {
		panic("failed to play chosen card")
	}
	trick.Plays = append(trick.Plays, Play{PlayerIndex: leaderIndex, Card: leadCard})
	for i := 1; i < len(r.Players); i++ {
		p := &r.Players[i]
		// TEMP: Randomly play one of valid hand cards (shuffled)
		card := p.LegalCards(&leadCard.Suit)[0]
		if !p.Play(card) {
			panic("failed to play chosen card")
		}
		trick.Plays = append(trick.Plays, Play{PlayerIndex: i, Card: card})
	}

	// Determine winning player, incrementing their trick count
	winner := trick.Winner()
	if winner >= 0 {
		r.Players[winner].Tricks++
	}
	// Log completed trick to history
	r.Tricks = append(r.Tricks, trick)

	fmt.Println(trick)
	fmt.Println(winner)

	// Calculate scores
	for i := range r.Players {
		p := &r.Players[i]
		if p.Bid == p.Tricks {
			p.Score += (10 + p.Bid)
		}
		// Currently no net-negative penalty, 0 points otherwise
		fmt.Println(p)
	}
}
