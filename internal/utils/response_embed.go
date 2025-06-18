package utils

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

func BuildEmbedMessageTemplate() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{Text: "© Clipify 2025",
			IconURL: "https://i.postimg.cc/B6GqPmgq/clipify2.png",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
