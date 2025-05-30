package handlers

import (
	"clipping-bot/internal/commands/account"
	"clipping-bot/internal/commands/campaign"
	"clipping-bot/internal/commands/stats"
	"clipping-bot/internal/commands/test"
	"clipping-bot/internal/commands/user"
	"clipping-bot/internal/commands/verification"
	"clipping-bot/internal/commands/video"
	"clipping-bot/internal/discord"
	"clipping-bot/internal/globalctx"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"strings"
)

func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, cancel := globalctx.ForRequest()
	defer cancel()

	// Handle buttons first
	if i.Type == discordgo.InteractionMessageComponent {
		if strings.HasPrefix(i.MessageComponentData().CustomID, "account_stats_") {
			stats.HandleAccountStatsButton(ctx, s, i)
			return
		}
		// ... handle other buttons here if needed
		return
	}

	// Handle slash commands
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		slog.Warn("failed to defer interaction", "error", err)
	}

	//exists, err := firebase.UserExists(ctx, i.Member.User.Username)
	//if err != nil {
	//	slog.Warn("failed to check if user exists", "error", err)
	//	discord.RespondToInteractionEmbedError(s, i, err.Error())
	//	return
	//}
	//if i.ApplicationCommandData().Name != "register" && !exists {
	//	slog.Error("user isn't registered")
	//	user.Register(ctx, s, i)
	//}

	switch i.ApplicationCommandData().Name {
	case "test":
		test.Test(s, i)
	case "register":
		user.Register(ctx, s, i)
	case "add-account":
		account.AddAccount(ctx, s, i)
	case "verify-account":
		verification.VerifyAccount(ctx, s, i)
	case "remove-verification":
		verification.RemoveVerification(ctx, s, i)
	case "add-video":
		video.AddVideo(ctx, s, i)
	case "stats":
		stats.GetStats(ctx, s, i)
	case "create-campaign":
		campaign.AddCampaign(ctx, s, i)
	default:
		discord.RespondToInteractionEmbedError(s, i, "Unknown command")
	}
}
