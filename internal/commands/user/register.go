package user

import (
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"clipping-bot/internal/utils"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func Register(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	ref, err := firebase.AddUser(ctx, i.Member.User.Username)
	if err != nil {
		slog.Error("Failed to add user", "error", err)
		discord.RespondToInteraction(s, i, utils.Capitalize(err.Error()))
		return
	}

	snapshot, err := ref.Get(ctx)
	if err != nil {
		slog.Error("Failed to retrieve user", "error", err)
		return
	}
	data := map[string]interface{}{
		"discord_username": snapshot.Data()["discord_username"],
	}

	discord.RespondToInteraction(s, i, fmt.Sprintf("Successfully registered user **%s**", data["discord_username"]))
}
