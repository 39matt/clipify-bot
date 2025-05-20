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

func RespondToInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate, responseMessage string) error {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseMessage,
		},
	})
	if err != nil {
		slog.Warn("failed to respond to interaction", "error", err)
		return err
	}
	return nil
}
