package game

import "slices"

type Game struct {
	Players []Player
	Rules   Rules
	Rounds  []Round
}

func NewGame(players []Player, rules Rules) *Game {
	return &Game{
		Players: slices.Clone(players),
		Rules:   rules,
	}
}
