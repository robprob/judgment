package game

import (
	"fmt"
	"io"
	"math/rand/v2"
	"slices"
)

type Game struct {
	Players []Player
	Rules   Rules
	Rounds  []*Round

	dealingRNG *rand.Rand

	Out io.Writer // Destination for game output, io.Discard by default
}

func NewGame(players []Player, rules Rules, seed uint64) *Game {
	return &Game{
		Players: slices.Clone(players),
		Rules:   rules,

		// Initialize seeded RNG
		dealingRNG: rand.New(rand.NewPCG(seed, 1)),

		Out: io.Discard,
	}
}

// Print final scores (in seat order) and state winner
func (g *Game) printFinalScores() {
	fmt.Fprintln(g.Out, "\nFinal scores")
	best := 0
	for i, p := range g.Players {
		fmt.Fprintf(g.Out, " %s: %d\n", p.Name, p.Score)
		if p.Score > g.Players[best].Score {
			best = i
		}
	}
	fmt.Fprintf(g.Out, "Winner: %s!\n", g.Players[best].Name)
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
		round.Out = g.Out

		fmt.Fprintf(g.Out, "\nRound %d: %d cards, trump %s\n",
			roundIndex+1, round.CardsPerPlayer, round.Trump)
		round.Run(g.dealingRNG)

		// Retain round players field as a shallow copy (frozen snapshot)
		round.Players = slices.Clone(round.Players)
		g.Rounds = append(g.Rounds, round)
	}

	g.printFinalScores()
}
