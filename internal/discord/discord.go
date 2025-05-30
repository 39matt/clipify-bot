package discord

import (
	"clipping-bot/internal/utils"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

//func RespondToInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate, responseMessage string) {
//	message := utils.Capitalize(responseMessage)
//	_, err := session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
//		Content: &message,
//	})
//	if err != nil {
//		slog.Warn("failed to respond to interaction", "error", err)
//	}
//}

func RespondToInteractionEmbed(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	title string, message string,
) {
	embed := utils.BuildEmbedMessageTemplate()
	embed.Title = title
	embed.Description = message
	embed.Color = 0x50C878
	_, err := session.InteractionResponseEdit(
		interaction.Interaction,
		&discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		},
	)
	if err != nil {
		slog.Warn("failed to respond to interaction", "error", err)
	}
}

func RespondToInteractionEmbedAndButtons(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	embed *discordgo.MessageEmbed,
	components []discordgo.MessageComponent,
) {
	embed.Color = 0x50C878
	_, err := session.InteractionResponseEdit(
		interaction.Interaction,
		&discordgo.WebhookEdit{
			Embeds:     &[]*discordgo.MessageEmbed{embed},
			Components: &components,
		},
	)
	if err != nil {
		slog.Warn("failed to respond to interaction", "error", err)
	}
}

func RespondToButtonInteractionEmbedAndButtons(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	embed *discordgo.MessageEmbed,
	components []discordgo.MessageComponent,
) {
	embed.Color = 0x50C878
	err := session.InteractionRespond(
		interaction.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{embed},
				Components: components,
			},
		},
	)
	if err != nil {
		slog.Warn("failed to respond to interaction", "error", err)
	}
}

func RespondToInteractionEmbedError(session *discordgo.Session, interaction *discordgo.InteractionCreate, errorMessage string) {
	embed := utils.BuildEmbedMessageTemplate()
	embed.Title = "❌ Error"
	embed.Description = utils.Capitalize(errorMessage)
	embed.Color = 0x50C878
	var components []discordgo.MessageComponent

	_, err := session.InteractionResponseEdit(
		interaction.Interaction,
		&discordgo.WebhookEdit{
			Embeds:     &[]*discordgo.MessageEmbed{embed},
			Components: &components,
		},
	)
	if err != nil {
		slog.Warn("failed to respond to interaction", "error", err)
	}
}
