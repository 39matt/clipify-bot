package discord

import (
	"clipping-bot/internal/config"
	"clipping-bot/internal/firebase"
	"clipping-bot/internal/globalctx"
	"clipping-bot/internal/models"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var (
	commandList = []*discordgo.ApplicationCommand{
		{
			Name:        "test",
			Description: "test",
		},
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
		//{
		//	Name:        "register",
		//	Description: "Register an account (linked to your discord account)",
		//},
		{
			Name:        "add-video",
			Description: "Track a new video",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "campaign",
					Description: "What campaign does your video belong to",
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
		{
			Name:        "create-campaign",
			Description: "Create a new campaign",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "influencer",
					Description: "Influencer's name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "activity",
					Description: "Activity description",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "budget",
					Description: "Budget for the campaign",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "per-million",
					Description: "Earnings per million views",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "max-submissions",
					Description: "Maximum number of submissions",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "max-earnings",
					Description: "Maximum total earnings",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "min-views-for-payout",
					Description: "Minimum views for payout",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "max-earnings-per-post",
					Description: "Maximum earnings per post",
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
	ctx, cancel := globalctx.ForRequest()
	defer cancel()

	for _, v := range commandList {
		cmd := *v

		if cmd.Name == "add-video" {
			campaigns, ids, getErr := firebase.GetCampaigns(ctx)
			if getErr != nil {
				slog.Error("failed to get campaigns", "error", getErr)
			}
			options := make([]*discordgo.ApplicationCommandOption, len(cmd.Options))
			copy(options, cmd.Options)
			options[0].Choices = nil
			for index, campaign := range campaigns {
				if index >= 25 {
					break
				}
				choice := &discordgo.ApplicationCommandOptionChoice{
					Name:  fmt.Sprintf("%s - %s", campaign.Influencer, campaign.Activity),
					Value: ids[index],
				}
				options[0].Choices = append(options[0].Choices, choice)
			}
			cmd.Options = options
		}

		_, err := Session.ApplicationCommandCreate(Session.State.User.ID, "", &cmd)
		if err != nil {
			fmt.Printf("Cannot create '%v' command: %v\n", cmd.Name, err)
		}
	}
}
