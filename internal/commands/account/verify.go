package account

import (
	"clipping-bot/internal/config"
	"clipping-bot/internal/models"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"net/http"
)

type TikTokUserResponse struct {
	UserInfo struct {
		User struct {
			Signature string `json:"signature"`
		} `json:"user"`
	} `json:"userInfo"`
}

func VerifyAccount(s *discordgo.Session, i *discordgo.InteractionCreate) {

	verificationInfo, exists := models.Verifications[i.Member.User.ID]
	if !exists {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("User **%s** doesn't have any pending verifications!", i.Member.User.Username),
			},
		})
		if err != nil {
			fmt.Println("Interaction response error:", err)
		}
		return
	}
	switch verificationInfo.Platform {
	case "TikTok":

		url := fmt.Sprintf("https://tiktok-api23.p.rapidapi.com/api/user/info?uniqueId=%s", models.Verifications[i.Member.User.ID].Username)

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("x-rapidapi-key", config.GetRapidApiKey())
		req.Header.Add("x-rapidapi-host", "tiktok-api23.p.rapidapi.com")

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		var result TikTokUserResponse
		err := json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("JSON unmarshal error:", err)
			return
		}

		bio := result.UserInfo.User.Signature
		fmt.Println("Bio:", bio)

		if models.Verifications[i.Member.User.ID].Code == bio {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Successfully verified **%s** (**%s**)!", verificationInfo.Username, verificationInfo.Platform),
				},
			})
			if err != nil {
				fmt.Println("Interaction response error:", err)
			}
			delete(models.Verifications, i.Member.User.ID)
		} else {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Please put **%s** on **%s** (**%s**) and then call '/verify-account' again", verificationInfo.Code, verificationInfo.Username, verificationInfo.Platform),
				},
			})
			if err != nil {
				fmt.Println("Interaction response error:", err)
			}
		}
	}

}
