package game

type Round struct {
	Players        []Player
	CardsPerPlayer int
	Trick          Trick
}

// Run a single round of Judgment
func (r *Round) Run() {
	// Instantiate a new deck
	// Shuffle and deal cards to each player
	// Players play through each trick
	// Calculate scores
}
