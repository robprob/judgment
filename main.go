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

// runBatch plays many silent-output games and prints each seat's average final score
func runBatch(games, numPlayers int, baseSeed uint64) {
	// Player totals log
	totals := make([]int, numPlayers)

	for i := range games {
		seed := baseSeed + uint64(i)
		players := newPlayers(numPlayers, seed)
		g := game.NewGame(players, game.DefaultRules(), seed)
		g.Run() // g.Out defaults to io.Discard
		for seat, p := range g.Players {
			totals[seat] += p.Score
		}
	}

	fmt.Printf("%d games (seeds %d to %d)\n", games, baseSeed, baseSeed+uint64(games-1))
	for seat, total := range totals {
		fmt.Printf("  Seat %d: average score %.2f\n", seat+1, float64(total)/float64(games))
	}
}

func main() {
	// CLI Validation
	seed := flag.Uint64("seed", 0, "RNG seed (0 = random)")
	games := flag.Int("games", 1, "number of games to simulate (more than 1 prints final summary only)")
	flag.Parse()

	if *seed == 0 {
		// Random seeding
		*seed = rand.Uint64()
	}

	numPlayers := 6

	if *games > 1 {
		runBatch(*games, numPlayers, *seed)
		return
	}

	fmt.Printf("Judgment time! (seed %d)\n", *seed)

	players := newPlayers(numPlayers, *seed)
	rules := game.DefaultRules()

	g := game.NewGame(players, rules, *seed)
	g.Out = os.Stdout
	g.Run()
}
