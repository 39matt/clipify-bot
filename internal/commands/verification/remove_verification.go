package verification

import (
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func RemoveVerification(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()

	err := firebase.RemoveVerification(ctx, i.Member.User.ID)

	if err != nil {
		respErr := discord.RespondToInteraction(s, i, err.Error())
		if respErr != nil {
			slog.Error("interaction respond failed", "error", respErr)
		}
		return
	}

	respErr := discord.RespondToInteraction(s, i, fmt.Sprintf("Successfully removed verification for **%s**", i.Member.User.Username))
	if respErr != nil {
		slog.Error("interaction respond failed", "error", respErr)
	}
}
