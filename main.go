package main

import (
	"clipping-bot/internal/bot"
	"clipping-bot/internal/firebase"
)

func main() {
	firebase.Initialize()
	bot.Start()
}
