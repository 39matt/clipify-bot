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
	var accountname string
	var platform models.Platform
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "platform":
			platform = models.Platform(option.StringValue())
		case "account-name":
			accountname = option.StringValue()
		}
	}

	accountAdded, accountErr := firebase.IsAccountAlreadyAdded(ctx, i.Member.User.Username, accountname, platform)
	if accountErr != nil {
		slog.Error("Account check failed", accountErr)
		discord.RespondToInteractionEmbedError(s, i, accountErr.Error())
	}

	accountExists, existsErr := firebase.IsAccountExists(ctx, platform, accountname)
	if existsErr != nil {
		slog.Error("Account check failed", accountErr)
		discord.RespondToInteractionEmbedError(s, i, accountErr.Error())
	}
	
	if accountAdded || accountExists {
		discord.RespondToInteractionEmbed(s, i, "⚠️ Warning", fmt.Sprintf("Account **%s** (**%s**) is already in use", accountname, platform))
		return
	}

	can, err := firebase.CanUserAddAccount(ctx, i.Member.User.Username, platform)
	if err != nil {
		slog.Error("can user add account function failed", "error", err)
		return
	}
	if !can {
		slog.Error("user already has 3 accounts", "error", err)
		discord.RespondToInteractionEmbed(s, i, "⚠️ Warning", "You already have 3 TikTok accounts bound to your user! Please remove one if you want to add another one.")
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
