package bot

import (
	"clipping-bot/internal/config"
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"clipping-bot/internal/handlers"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Start() {
	config.Load()
	if config.IsAppEnvironment(config.AppEnvironmentTest) {
		fmt.Println("App environment is test, aborting startup")
		return
	}

	discord.InitSession()

	firebase.Initialize()
	defer firebase.Close()

	addHandlers()
	discord.InitConnection()
	discord.RegisterCommands()

	fmt.Println("Bot is running. Press Ctrl + C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func addHandlers() {
	discord.Session.AddHandler(handlers.InteractionCreateHandler)
}
