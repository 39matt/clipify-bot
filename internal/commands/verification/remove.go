package verification

import (
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func RemoveVerification(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := firebase.RemoveVerification(ctx, i.Member.User.Username)

	if err != nil {
		discord.RespondToInteractionEmbedError(s, i, err.Error())
		return
	}

	discord.RespondToInteractionEmbed(s, i, "✅ Success", fmt.Sprintf("Successfully removed verification for **%s**", i.Member.User.Username))
}
