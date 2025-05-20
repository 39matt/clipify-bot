package account

import (
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func Register(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {

	ref, err := firebase.AddUser(ctx, i.Member.User.Username)
	if err != nil {
		slog.Error("Failed to add user", "error", err)
		respErr := discord.RespondToInteraction(s, i, err.Error())
		if respErr != nil {
			slog.Error("Failed to respond to interaction", "error", respErr)
		}
		return
	}

	snapshot, err := ref.Get(context.Background())
	if err != nil {
		slog.Error("Failed to retrieve user", "error", err)
	}
	data := map[string]interface{}{
		"discord_username": snapshot.Data()["discord_username"],
	}

	respErr := discord.RespondToInteraction(s, i, fmt.Sprintf("Successfully registered user **%s**", data["discord_username"]))
	if respErr != nil {
		slog.Error("Failed to respond to interaction", "error", respErr)
		return
	}
}
