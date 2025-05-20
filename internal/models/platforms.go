package models

import "github.com/bwmarrin/discordgo"

type Platform string

const (
	PlatformTikTok    Platform = "TikTok"
	PlatformInstagram Platform = "Instagram"
)

var PlatformChoices = []*discordgo.ApplicationCommandOptionChoice{
	{
		Name:  "TikTok",
		Value: string(PlatformTikTok),
	},
}
