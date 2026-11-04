package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"

	"github.com/robprob/judgment/internal/game"
)

// newPlayers creates numPlayers players, each with its own RNG stream derived from a starting seed
func newPlayers(numPlayers int, seed uint64) []game.Player {
	players := make([]game.Player, numPlayers)
	for i := range numPlayers {
		players[i].Name = "Player " + strconv.Itoa(i+1)
		// Per-player PCG instances
		rng := rand.New(rand.NewPCG(seed, uint64(100+i)))
		// Set strategy as uniform random selection with a maxBid of 3
		players[i].Strategy = game.NewRandomStrategy(rng, 3)
	}
	return players
}

func main() {
	// CLI Validation
	seed := flag.Uint64("seed", 0, "RNG seed (0 = random)")
	flag.Parse()

	if *seed == 0 {
		// Random seeding
		*seed = rand.Uint64()
	}
	fmt.Printf("Judgment time! (seed %d)\n", *seed)

	// Instantiate players
	numPlayers := 6
	players := newPlayers(numPlayers, *seed)

	rules := game.DefaultRules()

	g := game.NewGame(players, rules, *seed)
	g.Out = os.Stdout
	g.Run()
}
