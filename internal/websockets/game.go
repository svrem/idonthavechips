package websockets

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/svrem/idonthavechips/internal/utils"
)

type Game struct {
	Id string

	StartingCash int

	Started bool

	players []*Client
}

func NewGame(startingCash int) *Game {
	id := utils.RandSeq(6)

	return &Game{
		Id: id,

		StartingCash: startingCash,

		Started: false,

		players: []*Client{},
	}
}

func (g *Game) AddPlayer(client *Client) {
	g.players = append(g.players, client)
}

func (g *Game) RemovePlayer(player *Client) {
	for i, p := range g.players {
		if p == player {
			g.players = append(g.players[:i], g.players[i+1:]...)
			break
		}
	}
}

func (g *Game) Start() {
	if g.Started {
		return
	}

	for _, player := range g.players {
		player.send <- []byte(`{"type": "game-start"}`)
	}

	g.Started = true
}

func (g *Game) PlaceBet(player *Client, amount int) {
	if player.cash < amount {
		return
	}
	if player.currentBet >= amount {
		return
	}

	player.currentBet = amount
}

func (g *Game) CalculateResults(winners []*Client) {
	winnerPot := 0
	for _, player := range winners {
		winnerPot += player.currentBet
	}

	loserPot := 0
	for _, player := range g.players {
		if slices.Contains(winners, player) {
			continue
		}

		loss := min(player.currentBet, winnerPot)
		loserPot += loss
		slog.Info(fmt.Sprintf("Player %s lost %d", player.name, loss))
		player.UpdateCash(-loss)
	}

	for _, player := range winners {
		percentage := float64(player.currentBet) / float64(winnerPot)
		winnings := int(float64(loserPot) * percentage)
		slog.Info(fmt.Sprintf("Player %s won %d", player.name, winnings))
		player.UpdateCash(winnings)
	}

	for _, player := range g.players {
		player.currentBet = 0
	}
}
