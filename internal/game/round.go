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
	// Reset round-specific fields for all players
	for i := range r.Players {
		r.Players[i].RoundReset()
	}

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

	// Players state bids
	// TEMP: Random int [0, min(3, CardsPerPlayer)] capped at CardsPerPlayer
	// TODO: Enforce bid "hook" rule (will be on/true by default)
	for offset := range r.Players {
		// Find current player via wrap-around offset indexing
		playerIndex := (leaderIndex + offset) % len(r.Players)
		r.Players[playerIndex].Bid = rand.IntN(min(3, r.CardsPerPlayer) + 1)
	}

	// Players play through each trick

	// DEBUG: Single trick
	trick := Trick{Trump: r.Trump}

	var leadSuit *Suit
	for offset := range r.Players {
		// Find current player via wrap-around offset indexing
		playerIndex := (leaderIndex + offset) % len(r.Players)
		p := &r.Players[playerIndex]

		// Set lead/base suit after initial play
		if len(trick.Plays) == 1 {
			leadSuit = &trick.Plays[0].Card.Suit
		}
		// TEMP: Randomly play one of valid hand cards
		legal := p.LegalCards(leadSuit)
		card := legal[rand.IntN(len(legal))]
		if !p.Play(card) {
			panic("failed to play chosen card")
		}
		// Log play to current trick
		trick.Plays = append(trick.Plays, Play{PlayerIndex: playerIndex, Card: card})
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
