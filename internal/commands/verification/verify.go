package verification

import (
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"clipping-bot/internal/models"
	"clipping-bot/internal/utils"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"strings"
)

func VerifyAccount(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	verificationSnapshot, err := firebase.GetVerificationByDiscordUsername(ctx, i.Member.User.Username)
	if err != nil {
		slog.Error("Failed to get verification from firestore", "error", err)
		return
	}

	if verificationSnapshot == nil {
		discord.RespondToInteraction(s, i, fmt.Sprintf("User **%s** doesn't have any pending verifications!", i.Member.User.Username))
		return
	}

	var unverifiedAccount models.UnverifiedAccount
	err = verificationSnapshot.DataTo(&unverifiedAccount)
	if err != nil {
		slog.Error("error getting verification data", "error", err)
	}

	switch unverifiedAccount.Platform {
	case "TikTok":

		bio, err := utils.GetTiktokUserBio(unverifiedAccount.Username)

		if strings.Contains(bio, unverifiedAccount.Code) {
			err = firebase.RemoveVerification(ctx, i.Member.User.Username)
			if err == nil {
				slog.Info("Verification removed from firestore successfully!")
			}

			_, err = firebase.AddVerifiedAccount(ctx, i.Member.User.Username, models.Account{
				Username: unverifiedAccount.Username,
				Platform: unverifiedAccount.Platform,
				Link:     fmt.Sprintf("https://www.tiktok.com/@%s", unverifiedAccount.Username),
				Videos:   nil,
			})

			if err != nil {
				slog.Error("error adding verified account", "error", err)
				discord.RespondToInteraction(s, i, utils.Capitalize(err.Error()))
				return
			}
			discord.RespondToInteraction(s, i, fmt.Sprintf("Successfully verified **%s** (**%s**)!", unverifiedAccount.Username, unverifiedAccount.Platform))
		} else {
			discord.RespondToInteraction(s, i, fmt.Sprintf("Please put **%s** on **%s** (**%s**) and then call '/verify-account' again", unverifiedAccount.Code, unverifiedAccount.Username, unverifiedAccount.Platform))
		}
	}

}
