package verification

import (
	"clipify-bot/internal/discord"
	"clipify-bot/internal/firebase"
	"clipify-bot/internal/models"
	"clipify-bot/internal/utils"
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
		discord.RespondToInteractionEmbedError(s, i, fmt.Sprintf("User **%s** doesn't have any pending verifications!", i.Member.User.Username))
		return
	}

	var unverifiedAccount models.Verification
	err = verificationSnapshot.DataTo(&unverifiedAccount)
	if err != nil {
		slog.Error("error getting verification data", "error", err)
	}

	var bio string
	var account models.Account
	account.Username = unverifiedAccount.Username
	switch unverifiedAccount.Platform {
	case models.PlatformTikTok:
		var bioErr error
		bio, bioErr = utils.GetTiktokAccountBio(unverifiedAccount.Username)
		if bioErr != nil {
			slog.Error("error getting tiktok account bio", "error", bioErr)
			discord.RespondToInteractionEmbedError(s, i, bioErr.Error())
		}
		account.Link = fmt.Sprintf("https://www.tiktok.com/@%s", unverifiedAccount.Username)
		account.Platform = models.PlatformTikTok

	case models.PlatformInstagram:
		var bioErr error
		bio, bioErr = utils.GetInstagramAccountBio(unverifiedAccount.Username)
		if bioErr != nil {
			slog.Error("error getting instagram account bio", "error", bioErr)
			discord.RespondToInteractionEmbedError(s, i, bioErr.Error())
		}
		account.Link = fmt.Sprintf("https://instagram.com/%s", unverifiedAccount.Username)
		account.Platform = models.PlatformInstagram
	}

	if strings.Contains(bio, unverifiedAccount.Code) {
		err = firebase.RemoveVerification(ctx, i.Member.User.Username)
		if err == nil {
			slog.Info("Verification removed from firestore successfully!")
		}

		_, err = firebase.AddAccount(ctx, i.Member.User.Username, account)

		if err != nil {
			slog.Error("error adding verified account", "error", err)
			discord.RespondToInteractionEmbedError(s, i, err.Error())
			return
		}
		discord.RespondToInteractionEmbed(s, i, "✅ Success", fmt.Sprintf("Successfully verified **%s** (**%s**)!", unverifiedAccount.Username, unverifiedAccount.Platform))
	} else {
		discord.RespondToInteractionEmbed(s, i, "⚠️ Warning", fmt.Sprintf("Please put **%s** on **%s** (**%s**) and then call '/verify-account' again", unverifiedAccount.Code, unverifiedAccount.Username, unverifiedAccount.Platform))
	}

}
