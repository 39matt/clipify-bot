package account

import (
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"clipping-bot/internal/utils"
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
		discord.RespondToInteraction(s, i, "You already have 3 TikTok accounts bound to your user! Please remove one if you want to add another one.")
		return
	}

	accountExists, accountErr := firebase.IsAccountExists(ctx, i.Member.User.Username, accountname, platform)
	if accountErr != nil {
		slog.Error("Account check failed", accountErr)
		discord.RespondToInteraction(s, i, utils.Capitalize(accountErr.Error()))

	}
	if accountExists {
		discord.RespondToInteraction(s, i, fmt.Sprintf("Account **%s** (**%s**) is already in use", accountname, platform))
		return
	}

	rand.Seed(time.Now().UnixNano())
	verificationCode := rand.Intn(900000) + 100000
	doc, err := firebase.AddVerification(ctx, i.Member.User.Username, accountname, platform, verificationCode)
	if err != nil {
		slog.Error("error adding verification", "error", err)
		discord.RespondToInteraction(s, i, utils.Capitalize(err.Error()))
	}
	snapshot, err := doc.Get(ctx)
	if err != nil {
		slog.Error("error getting verification", "error", err)
	}
	data := snapshot.Data()
	discord.RespondToInteraction(s, i, fmt.Sprintf("Please add **%s** to your **%s** %s account bio, then use `/verify-account` to complete verification.", data["code"], data["username"], data["platform"]))
}
