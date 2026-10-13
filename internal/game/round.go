package game

import (
	"fmt"
	"math/rand/v2"
)

type Round struct {
	Players        []Player
	Rules          Rules
	CardsPerPlayer int
	StartingLeader int
	Trump          Suit
	Tricks         []Trick
}

func NewRound(players []Player, rules Rules, cardsPerPlayer, startingLeader int, trump Suit) *Round {
	return &Round{
		Players:        players,
		Rules:          rules,
		CardsPerPlayer: cardsPerPlayer,
		StartingLeader: startingLeader,
		Trump:          trump,
	}
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

	// Players state bids
	sumBids := 0
	maxBid := min(3, r.CardsPerPlayer)
	for offset := range r.Players {
		// Find current player via wrap-around offset indexing
		playerIndex := (r.StartingLeader + offset) % len(r.Players)
		// TEMP: Random int [0, min(3, CardsPerPlayer)] capped at CardsPerPlayer
		bid := rand.IntN(maxBid + 1)
		if offset == (len(r.Players)-1) && r.Rules.HookRule && (sumBids+bid) == r.CardsPerPlayer {
			// Bid is invalid via hook rule, player must choose another non-negative bid
			// TEMP: Random 50/50 +/- 1 unless bid already 0
			if bid == 0 {
				bid++
			} else if bid == maxBid {
				bid--
			} else if rand.IntN(2) == 0 {
				bid--
			} else {
				bid++
			}
		}
		sumBids += bid
		r.Players[playerIndex].Bid = bid
	}

	currentLeader := r.StartingLeader

	// Players play through all tricks, exhausting their hands
	for range r.CardsPerPlayer {

		trick := Trick{Trump: r.Trump}
		var leadSuit *Suit

		// Each player takes their turn playing a card for this trick
		for offset := range r.Players {
			// Find current player via wrap-around offset indexing
			playerIndex := (currentLeader + offset) % len(r.Players)
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

		// Increment winner trick count and establish leader for next trick
		winner := trick.Winner()
		if winner < 0 {
			panic("completed trick has no winner")
		}
		r.Players[winner].Tricks++
		currentLeader = winner

		// Log completed trick to history
		r.Tricks = append(r.Tricks, trick)

		fmt.Println(trick)
		fmt.Println(winner)
	}

	// Calculate end-of-round scores
	for i := range r.Players {
		p := &r.Players[i]
		if p.Bid == p.Tricks {
			p.Score += (10 + p.Bid)
		}
		// Currently no net-negative penalty, 0 points otherwise
		fmt.Println(p)
	}
}
