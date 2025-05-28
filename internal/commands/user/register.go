package user

import (
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func Register(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, err := firebase.AddUser(ctx, i.Member.User.Username)
	if err != nil {
		slog.Error("Failed to add user", "error", err)
		discord.RespondToInteractionEmbedError(s, i, err.Error())
		return
	}

	discord.RespondToInteractionEmbed(s, i, "✅ Success", fmt.Sprintf("Successfully registered user **%s**", i.Member.User.Username))
}
