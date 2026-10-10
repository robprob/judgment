package main

import (
	"fmt"

	"github.com/robprob/judgment/internal/game"
)

func main() {
	fmt.Println("Judgment time!")

	// Instantiate players
	numPlayers := 6
	players := make([]game.Player, numPlayers)
	for i := range numPlayers {
		players[i].ID = i + 1
	}

	round := game.Round{
		Players:        players,
		CardsPerPlayer: 7,
	}
	round.Run()
}
