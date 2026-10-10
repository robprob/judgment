package main

import (
	"fmt"

	"github.com/robprob/judgment/internal/game"
)

func main() {
	fmt.Println("Judgment time!")
	deck := game.NewDeck()
	fmt.Println(deck)
	deck.Shuffle()
	fmt.Println(deck)
}
