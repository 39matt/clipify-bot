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
	"strconv"
	"strings"
)

func VerifyAccount(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	verificationSnapshot, err := firebase.GetVerificationByDiscordID(ctx, i.Member.User.ID)
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

	var verificationInfo models.PendingVerification
	err = verificationSnapshot.DataTo(&verificationInfo)
	if err != nil {
		slog.Error("error getting verification data", "error", err)
	}

	switch verificationInfo.Platform {
	case "TikTok":

		url := fmt.Sprintf("https://tiktok-api23.p.rapidapi.com/api/user/info?uniqueId=%s", verificationInfo.AccountName)

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

		if strings.Contains(bio, strconv.Itoa(verificationInfo.Code)) {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Successfully verified **%s** (**%s**)!", verificationInfo.AccountName, verificationInfo.Platform),
				},
			})
			if err != nil {
				slog.Error("interaction response error", "error", err)
			}

			err = firebase.RemoveVerification(ctx, i.Member.User.ID)
			if err == nil {
				slog.Info("Verification removed from firestore successfully!")
			}

			//TODO add verified user to firebase

		} else {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Please put **%d** on **%s** (**%s**) and then call '/verify-account' again", verificationInfo.Code, verificationInfo.AccountName, verificationInfo.Platform),
				},
			})
			if err != nil {
				slog.Error("interaction response error", "error", err)
			}
		}
	}

}
