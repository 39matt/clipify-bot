package models

import "github.com/bwmarrin/discordgo"

type Platform string

const (
	PlatformYouTube Platform = "youtube"
	PlatformTwitch  Platform = "twitch"
	PlatformTikTok  Platform = "tiktok"
)

var PlatformChoices = []*discordgo.ApplicationCommandOptionChoice{
	{
		Name:  "YouTube",
		Value: string(PlatformYouTube),
	},
	{
		Name:  "Twitch",
		Value: string(PlatformTwitch),
	},
	{
		Name:  "TikTok",
		Value: string(PlatformTikTok),
	},
}
