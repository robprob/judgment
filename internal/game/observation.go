package game

// Slices are copies, safe to modify
// Seat-indexed unless otherwise stated

type PlayerBid struct {
	PlayerIndex int
	Bid         int
}

type BidObservation struct {
	// Player context
	PlayerIndex int
	Hand        []Card

	// Game context
	NumPlayers int
	Scores     []int // Running total before this round

	// Round context
	CardsPerPlayer int
	Blind          bool
	Trump          Suit

	// Strategy context
	LegalBids    []int
	PreviousBids []PlayerBid
}

type PlayObservation struct {
	// Player context
	PlayerIndex int
	Hand        []Card

	// Game context
	NumPlayers int
	Scores     []int

	// Round context
	CardsPerPlayer int
	Trump          Suit
	Bids           []int
	CurrentTrick   []Play
	PreviousTricks []Trick
	TricksWon      []int

	// Strategy context
	LegalCards []Card
}
