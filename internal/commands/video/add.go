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
	var platform, videoLink, campaignId string
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "campaign":
			campaignId = option.StringValue()
		case "platform":
			platform = option.StringValue()
		case "video-link":
			videoLink = option.StringValue()
		}
	}
	videoInfo := models.Video{}
	videoInfo.CampaignId = campaignId
	switch platform {
	case "TikTok":
		username := strings.Split(videoLink, "/")[3][1:]
		if strings.Contains(username, "?") {
			username = strings.Split(username, "?")[0]
		}

		videoId := strings.Split(videoLink, "/")[5]
		if strings.Contains(videoId, "?") {
			videoId = strings.Split(videoId, "?")[0]
		}

		names, err := firebase.GetUserAccountNames(ctx, i.Member.User.Username)
		if err != nil {
			slog.Error("Failed to get user account names", "error", err)
		}
		for index, name := range names {
			if name == username {
				break
			}
			if index == len(names)-1 {
				slog.Error("User doesn't have video's author in his verified accounts", "error", err)
				discord.RespondToInteractionEmbed(s, i, "⚠️ Warning", fmt.Sprintf("[This](%s) isn't your video! Please **add** this account or use another video uploaded by accounts you have set up.", videoLink))
				return
			}
		}

		var videoErr error
		videoInfo, videoErr = utils.GetTiktokVideoInfo(videoId)

		if videoErr != nil {
			discord.RespondToInteractionEmbedError(s, i, videoErr.Error())
			return
		}
	}

	_, err := firebase.AddVideo(ctx, i.Member.User.Username, videoInfo)
	if err != nil {
		discord.RespondToInteractionEmbedError(s, i, err.Error())
		return
	}
	embed := utils.BuildEmbedMessageTemplate()
	embed.Title = "Video upload results"
	embed.Description = fmt.Sprintf("%s\nHas been uploaded successfully ✅", videoInfo.Name)
	if len(videoInfo.Name) > 30 {
		embed.Description = fmt.Sprintf("%s\nHas been uploaded successfully ✅", videoInfo.Name[:27]+"...")
	}
	embed.URL = videoLink
	var components []discordgo.MessageComponent
	discord.RespondToInteractionEmbedAndButtons(s, i, embed, components)
}
