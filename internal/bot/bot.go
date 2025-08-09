package bot

import (
	"clipify-bot/internal/config"
	"clipify-bot/internal/discord"
	"clipify-bot/internal/firebase"
	"clipify-bot/internal/handlers"
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
