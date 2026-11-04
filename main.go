package main

import (
	"flag"
	"fmt"
	"math"
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
	rules := game.DefaultRules() // TEMP: Pass as function parameter in future

	// Player totals log
	numRounds := rules.StartingCards
	roundTotals := make([][]float64, numRounds) // [round][seat]
	for r := range roundTotals {
		roundTotals[r] = make([]float64, numPlayers)
	}
	totals := make([]float64, numPlayers)
	// Sum of squares
	squares := make([]float64, numPlayers)

	for i := range games {
		seed := baseSeed + uint64(i)
		players := newPlayers(numPlayers, seed)
		g := game.NewGame(players, rules, seed)
		g.Run() // g.Out defaults to io.Discard

		// Accumulate points per round (change in each seat's running score)
		prev := make([]int, numPlayers)
		for r, round := range g.Rounds {
			for seat, p := range round.Players {
				roundTotals[r][seat] += float64(p.Score - prev[seat])
				prev[seat] = p.Score
			}
		}

		// Accumulate each seat's final score and square
		for seat, p := range g.Players {
			score := float64(p.Score)
			totals[seat] += score
			squares[seat] += (score * score)
		}
	}

	fmt.Printf("%d games (seeds %d to %d)\n", games, baseSeed, baseSeed+uint64(games-1))

	n := float64(games)
	for seat := range numPlayers {
		mu := totals[seat] / n
		variance := (squares[seat] - (n * mu * mu)) / (n - 1)
		stdErr := math.Sqrt(variance / n)
		fmt.Printf("  Seat %d: average score %.2f +/- %.2f\n", seat+1, mu, 1.96*stdErr) // 95% confidence interval
	}

	fmt.Println("\nAverage points per round, by seat")
	for r, seats := range roundTotals {
		fmt.Printf("  Round %d:", r+1)
		for _, total := range seats {
			fmt.Printf(" %6.2f", total/n)
		}
		fmt.Println()
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
