package game

import (
	"math/rand/v2"
	"slices"
)

// RandomStrategy chooses uniformly at random among both legal bids and cards
// Randomizes legal bids 0..maxBid (negative means no cap/max allowed)
type RandomStrategy struct {
	rng    *rand.Rand // unique player RNG source
	maxBid int
}

// Compile-time check that RandomStrategy satisfies Strategy interface
var _ Strategy = (*RandomStrategy)(nil)

func NewRandomStrategy(rng *rand.Rand, maxBid int) *RandomStrategy {
	return &RandomStrategy{rng: rng, maxBid: maxBid}
}

func (s *RandomStrategy) ChooseBid(obs BidObservation) int {
	choices := obs.LegalBids
	if s.maxBid >= 0 {
		lowestLegal := slices.Min(obs.LegalBids)
		choices = slices.DeleteFunc(choices, func(b int) bool { return b > s.maxBid })
		if len(choices) == 0 {
			// Edge case: Hook rule removed every bid under the cap, take lowest legal bid
			return lowestLegal
		}
	}
	return choices[s.rng.IntN(len(choices))]
}

func (s *RandomStrategy) ChooseCard(obs PlayObservation) Card {
	return obs.LegalCards[s.rng.IntN(len(obs.LegalCards))]
}
