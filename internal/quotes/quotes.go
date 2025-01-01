package quotes

import "math/rand"

var quotes = []string{
	"Every spin could be the one that changes your life.",
	"The biggest wins come to those who dare to keep going.",
	"Quitting guarantees one thing—you’ll never know what could’ve been.",
	"The next roll of the dice might just rewrite your story.",
	"Winners are just players who never stopped believing in their luck.",
	"Patience turns small bets into big victories.",
	"You miss 100% of the chances you don’t take—so why not take one more?",
	"Dreams are built on risks, and risks are the heart of the game.",
	"Fortune is like a coin flip—sometimes, you’ve just got to stay in the game to see it land your way.",
	"The jackpot doesn’t go to the quitter; it goes to the one who stays in the fight.",
	"Every bet is a chance to turn the next page of your story.",
	"Every spin is a new opportunity to win.",
	"Every bet is a chance to turn the next page of your story.",
	"90% of the game is staying in the game.",
	"90% of gamblers quit before they hit big. Don’t be one of them.",
}

func RandomQuote() string {
	return quotes[rand.Intn(len(quotes))]
}
