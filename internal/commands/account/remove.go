package account

import (
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"clipping-bot/internal/models"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func RemoveAccount(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	var accountName string
	var platform models.Platform
	for _, option := range i.ApplicationCommandData().Options {
		if option.Name == "account-name" {
			accountName = option.StringValue()
		}
		if option.Name == "platform" {
			platform = models.Platform(option.StringValue())
		}
	}

	err := firebase.RemoveAccount(ctx, i.Member.User.Username, accountName, platform)
	if err != nil {
		discord.RespondToInteractionEmbedError(s, i, err.Error())
		return
	}
	discord.RespondToInteractionEmbed(s, i, "✅ Success", fmt.Sprintf("Successfully removed account: **%s (%s)**", accountName, platform))
}
