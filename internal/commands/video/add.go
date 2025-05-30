package video

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

func AddVideo(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	var videoLink, campaignId string
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "campaign":
			campaignId = option.StringValue()
		case "video-link":
			videoLink = option.StringValue()
		}
	}

	var platform models.Platform
	if strings.Contains(videoLink, "tiktok") {
		platform = models.PlatformTikTok
	} else if strings.Contains(videoLink, "instagram") {
		platform = models.PlatformInstagram
	} else {
		slog.Error("Link is not TT or IG")
		discord.RespondToInteractionEmbed(s, i, "⚠️ Warning", "Link you provided is not a Tiktok or Instagram link")
	}

	videoInfo := models.Video{}

	var videoErr error
	if platform == models.PlatformTikTok {
		//username := strings.Split(videoLink, "/")[3][1:]
		//if strings.Contains(username, "?") {
		//	username = strings.Split(username, "?")[0]
		//}
		videoId := strings.Split(videoLink, "/")[5]
		if strings.Contains(videoId, "?") {
			videoId = strings.Split(videoId, "?")[0]
		}
		videoInfo, videoErr = utils.GetTiktokVideoInfo(videoId)
	} else {
		videoId := strings.Split(videoLink, "/")[4]
		videoInfo, videoErr = utils.GetInstagramVideoInfo(videoId)

	}
	if videoErr != nil {
		discord.RespondToInteractionEmbedError(s, i, videoErr.Error())
		return
	}

	checkAccountExists(ctx, s, i, videoInfo.Author, videoLink)
	videoInfo.CampaignId = campaignId

	_, err := firebase.AddVideo(ctx, i.Member.User.Username, videoInfo.Author, platform, videoInfo)
	if err != nil {
		discord.RespondToInteractionEmbedError(s, i, err.Error())
		return
	}
	embed := utils.BuildEmbedMessageTemplate()
	embed.Title = "Video upload results"
	embed.Description = "**Clip #1**\n✅ Success"
	embed.URL = videoLink
	var components []discordgo.MessageComponent
	discord.RespondToInteractionEmbedAndButtons(s, i, embed, components)
}

func checkAccountExists(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, username string, videoLink string) bool {
	accounts, err := firebase.GetUserAccounts(ctx, i.Member.User.Username)
	if err != nil {
		slog.Error("Failed to get user account names", "error", err)
	}

	for index, account := range accounts {
		if account.Username == username {
			return true
		}
		if index == len(accounts)-1 {
			slog.Error("User doesn't have video's author in his verified accounts", "error", err)
			discord.RespondToInteractionEmbed(s, i, "⚠️ Warning", fmt.Sprintf("[This](%s) isn't your video! Please **add** this account or use another video uploaded by accounts you have set up.", videoLink))
			return false
		}
	}
	return false
}
