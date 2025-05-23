package verification

import (
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"clipping-bot/internal/utils"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func RemoveVerification(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := firebase.RemoveVerification(ctx, i.Member.User.Username)

	if err != nil {
		discord.RespondToInteraction(s, i, utils.Capitalize(err.Error()))
		return
	}

	discord.RespondToInteraction(s, i, fmt.Sprintf("Successfully removed verification for **%s**", i.Member.User.Username))
}
