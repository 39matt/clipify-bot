package account

import (
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"clipping-bot/internal/models"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"math/rand"
	"time"
)

func AddAccount(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	var platform, accountname string
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "platform":
			platform = option.StringValue()
		case "account-name":
			accountname = option.StringValue()
		}
	}

	can, err := firebase.CanUserAddTikTokAccount(ctx, i.Member.User.Username)
	if err != nil {
		slog.Error("can user add account function failed", "error", err)
		return
	}
	if !can {
		slog.Error("user already has 3 accounts", "error", err)
		discord.RespondToInteractionEmbed(s, i, "⚠️ Warning", "You already have 3 TikTok accounts bound to your user! Please remove one if you want to add another one.")
		return
	}

	accountExists, accountErr := firebase.IsAccountExists(ctx, i.Member.User.Username, accountname, platform)
	if accountErr != nil {
		slog.Error("Account check failed", accountErr)
		discord.RespondToInteractionEmbedError(s, i, accountErr.Error())

	}
	if accountExists {
		discord.RespondToInteractionEmbed(s, i, "⚠️ Warning", fmt.Sprintf("Account **%s** (**%s**) is already in use", accountname, platform))
		return
	}

	rand.Seed(time.Now().UnixNano())
	verificationCode := rand.Intn(900000) + 100000

	doc, err := firebase.AddVerification(ctx, i.Member.User.Username, accountname, platform, verificationCode)
	if err != nil {
		slog.Error("error adding verification", "error", err)
		discord.RespondToInteractionEmbedError(s, i, err.Error())
	}

	verificationSnapshot, err := doc.Get(ctx)
	if err != nil {
		slog.Error("error getting verification", "error", err)
	}

	var verificationData models.Verification
	err = verificationSnapshot.DataTo(&verificationData)
	if err != nil {
		slog.Error("error converting verification data to verification", "error", err)
		discord.RespondToInteractionEmbedError(s, i, "an unexpected error occurred")
		return
	}

	discord.RespondToInteractionEmbed(s, i, "✅ Success", fmt.Sprintf("Please add **%s** to your **%s** %s account bio, then use  `/verify-account` to complete verification.", verificationData.Code, verificationData.Username, verificationData.Platform))
}
