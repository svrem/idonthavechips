package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/svrem/idonthavechips/internal/quotes"
	"github.com/svrem/idonthavechips/internal/websockets"
)

func main() {
	IDToGameSocket := make(map[string]*websockets.GameSocket)

	http.HandleFunc("/ws/game/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		host := r.URL.Query().Get("host")
		name := r.URL.Query().Get("name")

		if name == "" && host != "true" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		gameSocket, ok := IDToGameSocket[id]
		if !ok {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}

		websockets.ServeWs(gameSocket, w, r, name, host == "true")
	})

	http.HandleFunc("/join-session", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		gameID := query.Get("game-id")
		name := query.Get("name")

		if name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)

			return
		}

		_, ok := IDToGameSocket[gameID]
		if !ok {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, "/game/"+gameID+"?name="+name, http.StatusFound)
	}))

	http.HandleFunc("/new-session", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/new-session.html")
	})

	http.HandleFunc("/api/new-session", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		cashAmountStr := query.Get("cash-amount")

		cashAmount, err := strconv.Atoi(cashAmountStr)
		if err != nil {
			http.Error(w, "Invalid cash amount", http.StatusBadRequest)
			return
		}

		if cashAmount < 1 {
			http.Error(w, "Cash amount must be greater than 0", http.StatusBadRequest)
			return
		}

		if cashAmount > 10000 {
			http.Error(w, "Cash amount must be less than 10000", http.StatusBadRequest)
			return
		}

		game := websockets.NewGame(cashAmount)
		gameSocket := websockets.NewGameSocket(game)

		IDToGameSocket[game.Id] = gameSocket

		go gameSocket.Run()

		url, err := url.Parse("/game/" + game.Id)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		q := url.Query()
		q.Set("host", "true")

		url.RawQuery = q.Encode()

		http.Redirect(w, r, url.String(), http.StatusFound)
	}))

	http.HandleFunc("/api/check-session", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		gameID := query.Get("game-id")

		gameSocket, ok := IDToGameSocket[gameID]
		if !ok {
			w.Write([]byte("false"))
			return
		}
		if gameSocket.Closed || gameSocket.GetGame().Started {
			w.Write([]byte("false"))
			return
		}

		w.Write([]byte("true"))
	}))

	type HostData struct {
		GameID string
	}

	type JoinData struct {
		GameID       string
		Name         string
		StartingCash int
	}

	http.HandleFunc("/game/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		query := r.URL.Query()
		host := query.Get("host")

		gameSocket, ok := IDToGameSocket[id]

		if !ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		if gameSocket.Closed {
			delete(IDToGameSocket, id)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		if host == "true" {
			template := template.Must(template.ParseFiles("views/host.html"))
			template.Execute(w, HostData{GameID: r.PathValue("id")})
			return
		}

		template := template.Must(template.ParseFiles("views/game.html"))
		template.Execute(w, JoinData{GameID: r.PathValue("id"), Name: query.Get("name"), StartingCash: gameSocket.GetGame().StartingCash})
	})

	http.HandleFunc("/join/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/join.html")
	})

	type IndexData struct {
		Quote string
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("views/index.html"))
		template.Execute(w, IndexData{Quote: quotes.RandomQuote()})

		// http.ServeFile(w, r, "views/index.html")
	})

	// assets
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/favicon.ico")
	})

	slog.Info("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
