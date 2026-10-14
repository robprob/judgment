package game

import (
	"math/rand/v2"
	"slices"
)

type Deck []Card

// Instantiate a standard 52-card playing deck
func NewDeck() Deck {
	deck := make(Deck, 0, 52)

	for s := Spades; s <= Hearts; s++ {
		for r := Two; r <= Ace; r++ {
			deck = append(deck, Card{Suit: s, Rank: r})
		}
	}
	return deck
}

// Randomly shuffle deck cards (Fisher-Yates)
func (d Deck) Shuffle(rng *rand.Rand) {
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

// DrawN draws N cards, removing them from the deck
func (d *Deck) DrawN(n int) []Card {
	if n <= 0 {
		// Requested 0 cards
		return nil
	}
	if n >= len(*d) {
		// Requested amount exceeds current deck length, empty and return entire deck
		draw := slices.Clone(*d)
		*d = nil
		return draw
	}

	draw := (*d)[:n] // Draw first N cards
	*d = (*d)[n:]    // Update remaining deck cards

	return draw
}
