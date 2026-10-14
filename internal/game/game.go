package game

import (
	"math/rand/v2"
	"slices"
)

type Game struct {
	Players []Player
	Rules   Rules
	Rounds  []*Round

	dealingRNG  *rand.Rand
	decisionRNG *rand.Rand
}

func NewGame(players []Player, rules Rules, seed uint64) *Game {
	return &Game{
		Players: slices.Clone(players),
		Rules:   rules,

		// Initialize seeded RNG
		dealingRNG:  rand.New(rand.NewPCG(seed, 1)),
		decisionRNG: rand.New(rand.NewPCG(seed, 2)),
	}
}

// Run a full game of Judgment
func (g *Game) Run() {
	// Iterate through rounds of the game
	for roundIndex := range g.Rules.StartingCards {
		round := NewRound(
			g.Players,
			g.Rules,
			g.Rules.StartingCards-roundIndex,
			roundIndex%len(g.Players),
			g.Rules.TrumpOrder[roundIndex%len(g.Rules.TrumpOrder)],
		)
		round.Run(g.dealingRNG, g.decisionRNG)

		// Retain round players field as a shallow copy (frozen snapshot)
		round.Players = slices.Clone(round.Players)
		g.Rounds = append(g.Rounds, round)
	}
}
