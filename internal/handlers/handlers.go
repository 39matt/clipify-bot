package handlers

import (
	"clipping-bot/internal/commands/account"
	"clipping-bot/internal/commands/stats"
	"clipping-bot/internal/commands/user"
	"clipping-bot/internal/commands/verification"
	"clipping-bot/internal/commands/video"
	"clipping-bot/internal/globalctx"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"strings"
)

func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, cancel := globalctx.ForRequest()
	defer cancel()

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	if i.Type == discordgo.InteractionMessageComponent {
		if strings.HasPrefix(i.MessageComponentData().CustomID, "account_stats_") {
			stats.HandleAccountStatsButton(ctx, s, i)
			return
		}
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		slog.Warn("failed to defer interaction", "error", err)
	}

	switch i.ApplicationCommandData().Name {
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

	}
}
