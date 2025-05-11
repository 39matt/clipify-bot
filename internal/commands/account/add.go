package account

import (
	"clipping-bot/internal/models"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
)

func AddAccount(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var platform, accountname string
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "platform":
			platform = option.StringValue()
		case "accountname":
			accountname = option.StringValue()
		}
	}

	_, exists := models.Verifications[i.Member.User.ID]
	if exists {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("You already have one pending verification! Please finish the verification with **%s** on **%s**", models.Verifications[i.Member.User.ID].Username, models.Verifications[i.Member.User.ID].Platform),
			},
		})
	}

	verificationCode := generateRandomNumber()
	models.Verifications[i.Member.User.ID] = models.PendingVerification{
		Code:     verificationCode,
		Platform: platform,
		Username: accountname,
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Please add **%s** to your %s account bio, then use `/verify-account` to complete verification.", verificationCode, platform),
		},
	})
}

func generateRandomNumber() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
