package utils

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

func BuildEmbedMessageTemplate() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{Text: "© Clipify 2025",
			IconURL: "https://ibb.co/wNsJm4gj",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
