package discord

import (
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

//func SendChannelMessage(channelID string, message string) {
//	_, err := Session.ChannelMessageSend(channelID, message)
//	if err != nil {
//		slog.Warn("failed to send message to channel", "channelId", channelID, "message", message, "error", err)
//	}
//}

func RespondToInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate, responseMessage string) {
	_, err := session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Content: &responseMessage,
	})
	if err != nil {
		slog.Warn("failed to respond to interaction", "error", err)
	}
}
