# Judgment

A Go implementation of **Judgment**, a trick-taking card game centered around predicting exactly how many tricks you'll win each round.

A personal project exploring game simulation, probability, and eventually Monte Carlo analysis of different playing strategies.

## Rules

Based on the variation played with friends, at least that I remember :)

- Standard 52-card deck.
- Six players, starting with seven cards each and decreasing each round until a final, single-card round.
- Players bid on the number of tricks they expect to win.
- Players must play the lead suit when possible; otherwise, they may play any card.
- Trump cards beat non-trump cards, with the highest trump winning the trick.
- Successfully matching a bid awards **10 + bid** points; otherwise, zero.
- The player with the highest cumulative score across rounds wins.

Trump suit rotates between rounds. The final bidder is restricted from making the sum of bids equal the number of available tricks (the *hook rule*).

## Running

```bash
go run .
```

The current program runs an experimental single round and prints player hands, completed tricks, and scores.

## Project Status

Early development. Currently focused on building a simple, correct game engine before expanding into configurable gameplay and strategy simulation.