package game

// Strategy decides bidding and played card decisions for a player
// Implementations must return an element from the set of LegaBids / LegalCards
// Engine validates chosen values, panics otherwise
type Strategy interface {
	ChooseBid(obs BidObservation) int
	ChooseCard(obs PlayObservation) Card
}
