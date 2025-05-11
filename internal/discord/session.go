package discord

import (
	"clipping-bot/internal/models"
	"fmt"
	"log/slog"

	"clipping-bot/internal/config"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "add-account",
			Description: "Add a new account",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "platform",
					Description: "Platform of the account (e.g. YouTube, Twitch)",
					Choices:     models.PlatformChoices,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "accountname",
					Description: "Name of the account",
					Required:    true,
				},
			},
		},
	}
)

var Session *discordgo.Session

func InitSession() {
	var err error
	Session, err = discordgo.New("Bot " + config.GetDiscordToken()) // Initializing discord session
	if err != nil {
		slog.Error("failed to create discord session", "error", err)
	}

	Session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuildMessageTyping | discordgo.IntentGuildVoiceStates | discordgo.IntentGuilds
}

func InitConnection() {
	if err := Session.Open(); err != nil { // Creating a connection
		slog.Error("failed to create websocket connection to discord", "error", err)
		return
	}
}

func RegisterCommands() {
	for _, v := range commands {
		_, err := Session.ApplicationCommandCreate(Session.State.User.ID, "", v) // "" means global command
		if err != nil {
			fmt.Printf("Cannot create '%v' command: %v\n", v.Name, err)
		}
	}
}
