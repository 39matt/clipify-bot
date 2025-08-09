package campaign

import (
	"clipify-bot/internal/discord"
	"clipify-bot/internal/firebase"
	"clipify-bot/internal/models"
	"clipify-bot/internal/utils"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

func AddCampaign(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	var campaign models.Campaign
	campaign.CreatedAt = time.Now().Format(time.RFC3339)
	campaign.Progress = 0
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "influencer":
			campaign.Influencer = option.StringValue()
		case "activity":
			campaign.Activity = option.StringValue()
		case "budget":
			campaign.Budget = option.StringValue()
		case "per-million":
			campaign.PerMillion = option.FloatValue()
		case "max-submissions":
			campaign.MaxSubmissions = option.FloatValue()
		case "max-earnings":
			campaign.MaxEarnings = option.FloatValue()
		case "min-views-for-payout":
			campaign.MinViewsForPayout = option.FloatValue()
		case "max-earnings-per-post":
			campaign.MaxEarningsPerPost = option.FloatValue()
		}
	}

	_, err := firebase.AddCampaign(ctx, campaign)
	if err != nil {
		embed := utils.BuildEmbedMessageTemplate()
		embed.Title = "Error"
		embed.Description = err.Error()
		discord.RespondToInteractionEmbedAndButtons(s, i, embed, []discordgo.MessageComponent{})
	}
	embed := utils.BuildEmbedMessageTemplate()
	embed.Title = "Successfully added campaign ✅"
	embed.Description = fmt.Sprintf(
		"Influencer: %s\nActivity: %s\nBudget: %s\nProgress: %.2f%%\nPer Million: %.2f\nMax Submissions: %.0f\nMax Earnings: %.2f\nMax Earnings Per Post: %.2f\nMin Views For Payout: %.0f",
		campaign.Influencer,
		campaign.Activity,
		campaign.Budget,
		campaign.Progress,
		campaign.PerMillion,
		campaign.MaxSubmissions,
		campaign.MaxEarnings,
		campaign.MaxEarningsPerPost,
		campaign.MinViewsForPayout,
	)

	discord.RespondToInteractionEmbedAndButtons(s, i, embed, []discordgo.MessageComponent{})
}
