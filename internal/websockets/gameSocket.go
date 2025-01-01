package websockets

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type GameSocket struct {
	host    *Client
	players map[*Client]bool

	game *Game

	Closed bool

	register      chan *Client
	unregister    chan *Client
	makeHost      chan *Client
	handleMessage chan *Message
}

func (gs *GameSocket) GetGame() *Game {
	return gs.game
}

func NewGameSocket(game *Game) *GameSocket {
	gameSocket := GameSocket{
		host:    nil,
		players: make(map[*Client]bool),

		game: game,

		Closed: false,

		register:      make(chan *Client),
		unregister:    make(chan *Client),
		makeHost:      make(chan *Client),
		handleMessage: make(chan *Message),
	}

	return &gameSocket
}

type PlayerDataMessage struct {
	Type string       `json:"type"`
	Data []PlayerData `json:"data"`
}

type PlayerData struct {
	Name string `json:"name"`
	Cash int    `json:"cash"`
}

func (gs *GameSocket) PlayerData() ([]byte, error) {
	var data []PlayerData

	for _, player := range gs.game.players {
		data = append(data, PlayerData{
			Name: player.name,
			Cash: player.cash,
		})
	}

	playerDataMessage := PlayerDataMessage{
		Type: "player-list-update",
		Data: data,
	}

	return json.Marshal(playerDataMessage)
}

func (gs *GameSocket) Run() {
	defer func() {
		gs.Closed = true
		for client := range gs.players {
			delete(gs.players, client)
			close(client.send)
		}
	}()

	for {
		select {
		case client := <-gs.register:
			gs.players[client] = true

			gs.game.AddPlayer(client)

			if gs.host != nil {
				data, err := gs.PlayerData()

				if err != nil {
					slog.Error("Error marshalling player data")
					return
				}

				gs.host.send <- data
			}

		case client := <-gs.unregister:
			if _, ok := gs.players[client]; ok {
				gs.game.RemovePlayer(client)

				if gs.host != nil {
					data, err := gs.PlayerData()

					if err != nil {
						slog.Error("Error marshalling player data")
						return
					}

					gs.host.send <- data
				}

				delete(gs.players, client)
				close(client.send)

				if gs.host == client {
					gs.host = nil
					return
				}
			}

			// if gs.host == client {
			// 	gs.host = nil
			// }

			// for client := range gs.players {
			// 	close(client.send)
			// }

			// return
		case client := <-gs.makeHost:
			if gs.host == nil {
				gs.host = client
				gs.game.RemovePlayer(client)
			}

		case message := <-gs.handleMessage:
			fmt.Println(string(message.data))

			var data interface{}
			if err := json.Unmarshal(message.data, &data); err != nil {
				fmt.Println(err)
				slog.Error("Error unmarshalling message data")
				continue
			}

			dataMap, ok := data.(map[string]interface{})
			if !ok {
				slog.Error("Error asserting message data to map")
				continue
			}

			switch dataMap["type"] {
			case "start":
				if gs.host == message.client {
					gs.game.Start()
					gs.host.send <- []byte(`{"type": "game-start"}`)
				}
			case "place-bet":
				slog.Info("Place bet " + message.client.name)
				bet := int(dataMap["bet"].(float64))

				if message.client.cash < bet {
					slog.Info("Not enough cash")
					continue
				}

				if message.client.currentBet >= bet {
					slog.Info("Bet is less than current bet")
					continue
				}

				message.client.currentBet = bet

				gs.host.send <- []byte(fmt.Sprintf(`{"type": "bet-placed", "name": "%s", "bet": %d, "id": "%s"}`, message.client.name, bet, message.client.id))
			case "declare-winners":
				if gs.host != message.client {
					continue
				}

				winners := dataMap["winners"].([]interface{})

				var winnerClients []*Client

				for _, winner := range winners {
					for client := range gs.players {
						if client.id == winner {
							winnerClients = append(winnerClients, client)
						}
					}
				}

				gs.game.CalculateResults(winnerClients)

				gs.host.send <- []byte(`{"type": "game-end"}`)
			}
		}

	}
}
