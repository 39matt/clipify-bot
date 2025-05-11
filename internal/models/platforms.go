package models

import "github.com/bwmarrin/discordgo"

type Platform string

const (
	PlatformTikTok  Platform = "TikTok"
	PlatformYouTube Platform = "YouTube"
)

var PlatformChoices = []*discordgo.ApplicationCommandOptionChoice{
	{
		Name:  "TikTok",
		Value: string(PlatformTikTok),
	},
}
