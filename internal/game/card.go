package game

import "fmt"

// Defines standard 52-card deck ordered suits and ranks

type Suit uint8

const (
	Spades Suit = iota
	Diamonds
	Clubs
	Hearts
)

func (s Suit) String() string {
	switch s {
	case Spades:
		return "♠"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	case Hearts:
		return "♥"
	default:
		return fmt.Sprintf("Unknown Suit (%d)", s)
	}
}

type Rank uint8

const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func (r Rank) String() string {
	switch r {
	case Two:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Nine:
		return "9"
	case Ten:
		return "10"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	default:
		return fmt.Sprintf("Unknown Rank (%d)", r)
	}
}

// Define card structure holding suit/rank

type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Suit, c.Rank)
}
