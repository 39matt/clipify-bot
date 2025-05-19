package account

import (
	"clipping-bot/internal/firebase"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"math/rand"
)

func AddAccount(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	var platform, accountname string
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "platform":
			platform = option.StringValue()
		case "accountname":
			accountname = option.StringValue()
		}
	}

	verificationCode := rand.Intn(900000) + 100000
	doc, err := firebase.AddVerification(ctx, i.Member.User.ID, accountname, platform, verificationCode)
	if err != nil {
		slog.Error("error adding verification", "error", err)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
			},
		})
		return
	}
	snapshot, err := doc.Get(ctx)
	if err != nil {
		slog.Error("error getting verification", "error", err)
	}
	data := snapshot.Data()
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Please add **%d** to your **%s** %s account bio, then use `/verify-account` to complete verification.", data["code"], data["accountname"], data["platform"]),
		},
	})
}
