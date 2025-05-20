package verification

import (
	"clipping-bot/internal/config"
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"clipping-bot/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

func VerifyAccount(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	verificationSnapshot, err := firebase.GetVerificationByDiscordUsername(ctx, i.Member.User.Username)
	if err != nil {
		slog.Error("Failed to get verification from firestore", "error", err)
		return
	}

	if verificationSnapshot == nil {
		respErr := discord.RespondToInteraction(s, i, fmt.Sprintf("User **%s** doesn't have any pending verifications!", i.Member.User.Username))
		if respErr != nil {
			slog.Error("interaction respond failed", "error", respErr)
		}
		return
	}

	var unverifiedAccount models.UnverifiedAccount
	err = verificationSnapshot.DataTo(&unverifiedAccount)
	if err != nil {
		slog.Error("error getting verification data", "error", err)
	}

	switch unverifiedAccount.Platform {
	case "TikTok":

		url := fmt.Sprintf("https://tiktok-api23.p.rapidapi.com/api/user/info?uniqueId=%s", unverifiedAccount.Username)

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("x-rapidapi-key", config.GetRapidApiKey())
		req.Header.Add("x-rapidapi-host", "tiktok-api23.p.rapidapi.com")

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)

		var result models.TikTokUserResponse
		jsonErr := json.Unmarshal(body, &result)
		if jsonErr != nil {
			slog.Error("JSON unmarshal error", "error", jsonErr)
		}

		bio := result.UserInfo.User.Signature

		if strings.Contains(bio, unverifiedAccount.Code) {
			err = discord.RespondToInteraction(s, i, fmt.Sprintf("Successfully verified **%s** (**%s**)!", unverifiedAccount.Username, unverifiedAccount.Platform))
			if err != nil {
				slog.Error("interaction response error", "error", err)
			}

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
			}

		} else {
			err = discord.RespondToInteraction(s, i, fmt.Sprintf("Please put **%s** on **%s** (**%s**) and then call '/verify-account' again", unverifiedAccount.Code, unverifiedAccount.Username, unverifiedAccount.Platform))
			if err != nil {
				slog.Error("interaction response error", "error", err)
			}
		}
	}

}
