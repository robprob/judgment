package game

import (
	"fmt"
	"io"
	"math/rand/v2"
	"slices"
)

type Round struct {
	Players        []Player
	Rules          Rules
	CardsPerPlayer int
	StartingLeader int
	Trump          Suit
	Tricks         []Trick

	Out io.Writer // Destination for round output, io.Discard by default
}

func NewRound(players []Player, rules Rules, cardsPerPlayer, startingLeader int, trump Suit) *Round {
	return &Round{
		Players:        players,
		Rules:          rules,
		CardsPerPlayer: cardsPerPlayer,
		StartingLeader: startingLeader,
		Trump:          trump,
		Out:            io.Discard,
	}
}

// Run a single round of Judgment
func (r *Round) Run(dealingRNG *rand.Rand) {
	// Reset round-specific fields for all players
	for i := range r.Players {
		r.Players[i].RoundReset()
	}

	// State/Config
	numPlayers := len(r.Players)
	blind := r.CardsPerPlayer == 1 && r.Rules.BlindOneCardRound
	// Parse current score for all players
	scores := make([]int, numPlayers)
	for playerIndex := range r.Players {
		scores[playerIndex] = r.Players[playerIndex].Score
	}

	// Instantiate a new deck
	deck := NewDeck()
	// Shuffle and deal cards to each player
	deck.Shuffle(dealingRNG)
	for i := range r.Players {
		r.Players[i].Hand = deck.DrawN(r.CardsPerPlayer)
	}

	// Players state bids
	var playerBids []PlayerBid
	sumSoFar := 0
	for offset := range r.Players {
		// Find current player via wrap-around offset indexing
		playerIndex := (r.StartingLeader + offset) % numPlayers
		p := &r.Players[playerIndex]

		// Enumerate legal bids for this player
		isLast := offset == (numPlayers - 1)
		legal := legalBids(r.CardsPerPlayer, sumSoFar, isLast, r.Rules.HookRule)

		// Build bid observation context available to this player
		var hand []Card
		if !blind {
			hand = slices.Clone(p.Hand)
		}
		bidObs := BidObservation{
			// Player context
			PlayerIndex: playerIndex,
			Hand:        hand,

			// Game context
			NumPlayers: numPlayers,
			Scores:     slices.Clone(scores),

			// Round context
			CardsPerPlayer: r.CardsPerPlayer,
			Blind:          blind,
			Trump:          r.Trump,

			// Strategy context
			LegalBids:    slices.Clone(legal),
			PreviousBids: slices.Clone(playerBids),
		}

		// Player chooses a bid based on their strategy
		bid := p.Strategy.ChooseBid(bidObs)
		// Validate chosen bid
		if !slices.Contains(legal, bid) {
			panic(fmt.Sprintf("%s chose illegal bid: %d", p.Name, bid))
		}

		playerBids = append(playerBids, PlayerBid{playerIndex, bid})
		sumSoFar += bid
		p.Bid = bid
	}

	currentLeader := r.StartingLeader

	// Players play through all tricks, exhausting their hands
	for range r.CardsPerPlayer {

		trick := Trick{Trump: r.Trump}
		var leadSuit *Suit

		// Each player takes their turn playing a card for this trick
		for offset := range r.Players {
			// Find current player via wrap-around offset indexing
			playerIndex := (currentLeader + offset) % numPlayers
			p := &r.Players[playerIndex]

			// Set lead/base suit after initial play
			if len(trick.Plays) == 1 {
				leadSuit = &trick.Plays[0].Card.Suit
			}

			// Enumerate legal cards for this player
			legal := p.LegalCards(leadSuit)

			// Build play observation context available to this player
			playObs := PlayObservation{
				// Player context
				PlayerIndex: playerIndex,
				Hand:        slices.Clone(p.Hand),

				// Game context
				NumPlayers: numPlayers,
				Scores:     slices.Clone(scores),

				// Round context
				CardsPerPlayer: r.CardsPerPlayer,
				Trump:          r.Trump,
				Bids:           r.bids(),
				CurrentTrick:   trick.clone(),
				PreviousTricks: r.previousTricks(),
				TricksWon:      r.tricksWon(),

				// Strategy context
				LegalCards: slices.Clone(legal),
			}

			// Player chooses a card based on their strategy
			card := p.Strategy.ChooseCard(playObs)
			if !slices.Contains(legal, card) {
				panic(fmt.Sprintf("%s chose illegal card: %s", p.Name, card))
			}
			if !p.Play(card) {
				panic(fmt.Sprintf("%s chose card not in hand: %s", p.Name, card))
			}
			// Log play to current trick
			trick.Plays = append(trick.Plays, Play{PlayerIndex: playerIndex, Card: card})
		}

		// Increment winner trick count and establish leader for next trick
		winner := trick.Winner()
		if winner < 0 {
			panic("completed trick has no winner")
		}
		r.Players[winner].TricksWon++

		currentLeader = winner // winning player leads off next trick

		// Log completed trick to history
		r.Tricks = append(r.Tricks, trick)
	}

	// Calculate end-of-round scores
	for i := range r.Players {
		p := &r.Players[i]
		points := 0
		if p.Bid == p.TricksWon {
			points = 10 + p.Bid
		}
		// Currently no net-negative penalty, 0 points otherwise
		p.Score += points
		fmt.Fprintf(r.Out, " %s: bid %d, won %d, +%d (total %d)\n",
			p.Name, p.Bid, p.TricksWon, points, p.Score)
	}
}

// legalBids returns the bids currently available to a player (in ascending order)
// Hook Rule applies, if enabled
func legalBids(cardsPerPlayer, sumSoFar int, isLast, hookRule bool) []int {
	bids := make([]int, 0, cardsPerPlayer+1)
	for b := range cardsPerPlayer + 1 {
		if hookRule && isLast && sumSoFar+b == cardsPerPlayer {
			continue // Forbidden by hook rule
		}
		bids = append(bids, b)
	}
	return bids
}

// bids returns a fresh copy of each player's chosen bid for this round
func (r *Round) bids() []int {
	bids := make([]int, len(r.Players))
	for i, p := range r.Players {
		bids[i] = p.Bid
	}
	return bids
}

// previousTricks returns a deep copy of this round's completed tricks
func (r *Round) previousTricks() []Trick {
	tricks := make([]Trick, len(r.Tricks))
	for i, t := range r.Tricks {
		tricks[i] = t.clone()
	}
	return tricks
}

// tricksWon returns a fresh copy of each player's tricks won so far
// Scopes an intentionally reduced amount of per-player information
func (r *Round) tricksWon() []int {
	won := make([]int, len(r.Players))
	for i, p := range r.Players {
		won[i] = p.TricksWon
	}
	return won
}
