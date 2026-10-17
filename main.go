package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/robprob/judgment/internal/game"
)

func main() {
	// CLI Validation
	seed := flag.Uint64("seed", 0, "RNG seed (0 = random)")
	flag.Parse()

	if *seed == 0 {
		// Random seeding
		*seed = rand.Uint64()
	}
	fmt.Printf("Seed: %d\n", *seed)

	fmt.Println("Judgment time!")

	// Instantiate players
	numPlayers := 6
	players := make([]game.Player, numPlayers)
	for i := range numPlayers {
		players[i].Name = "Player " + strconv.Itoa(i+1)
	}

	rules := game.DefaultRules()

	game := game.NewGame(players, rules, *seed)
	game.Run()
}
