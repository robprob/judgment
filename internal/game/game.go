package game

import "slices"

type Game struct {
	Players []Player
	Rules   Rules
	Rounds  []*Round
}

func NewGame(players []Player, rules Rules) *Game {
	return &Game{
		Players: slices.Clone(players),
		Rules:   rules,
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
		round.Run()

		// Retain round players field as a shallow copy (frozen snapshot)
		round.Players = slices.Clone(round.Players)
		g.Rounds = append(g.Rounds, round)
	}
}
