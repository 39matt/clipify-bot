package test

import (
	"clipping-bot/internal/discord"
	"github.com/bwmarrin/discordgo"
)

func Test(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{Title: "Successfully added video", Description: "video info name"}
	var components []discordgo.MessageComponent
	discord.RespondToInteractionEmbedAndButtons(s, i, embed, components)
}
