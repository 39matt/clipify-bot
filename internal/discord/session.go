package discord

import (
	"clipping-bot/internal/config"
	"clipping-bot/internal/models"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var (
	commandList = []*discordgo.ApplicationCommand{
		{
			Name:        "add-account",
			Description: "Add a new account",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "platform",
					Description: "Platform of the account (e.g. TikTok, Instagram)",
					Choices:     models.PlatformChoices,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "account-name",
					Description: "Name of the account",
					Required:    true,
				},
			},
		},
		{
			Name:        "verify-account",
			Description: "Check if bio has the verification code",
		},
		{
			Name:        "remove-verification",
			Description: "Remove active verification",
		},
		{
			Name:        "register",
			Description: "Register an account (linked to your discord account)",
		},
		{
			Name:        "add-video",
			Description: "Track a new video",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "platform",
					Description: "Platform of the account (e.g. ikTok, Instagram)",
					Choices:     models.PlatformChoices,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "video-link",
					Description: "Link to the video you uploaded",
					Required:    true,
				},
			},
		},
		{
			Name:        "stats",
			Description: "Get your stats",
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
	// Register new commands
	for _, v := range commandList {
		_, err := Session.ApplicationCommandCreate(Session.State.User.ID, "", v)
		if err != nil {
			fmt.Printf("Cannot create '%v' command: %v\n", v.Name, err)
		}
	}
}
