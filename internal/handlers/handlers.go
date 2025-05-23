package handlers

import (
	"clipping-bot/internal/commands/account"
	"clipping-bot/internal/commands/user"
	"clipping-bot/internal/commands/verification"
	"clipping-bot/internal/commands/video"
	"clipping-bot/internal/globalctx"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		slog.Warn("failed to defer interaction", "error", err)
	}

	ctx, cancel := globalctx.ForRequest()
	defer cancel()

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
	}
}
