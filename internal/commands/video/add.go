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
	var platform, videoLink string
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "platform":
			platform = option.StringValue()
		case "video-link":
			videoLink = option.StringValue()
		}
	}
	videoInfo := models.Video{}
	switch platform {
	case "TikTok":
		username := strings.Split(videoLink, "/")[3][1:]
		videoId := strings.Split(videoLink, "/")[5]

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
				discord.RespondToInteraction(s, i, fmt.Sprintf("[This](%s) isn't your video! Please **add** this account or use another video uploaded by accounts you have set up.", videoLink))
				return
			}
		}

		var videoErr error
		videoInfo, videoErr = utils.GetTiktokVideoInfo(videoId)

		if videoErr != nil {
			discord.RespondToInteraction(s, i, utils.Capitalize(videoErr.Error()))
			return
		}
	}

	_, err := firebase.AddVideo(ctx, i.Member.User.Username, videoInfo)
	if err != nil {
		discord.RespondToInteraction(s, i, utils.Capitalize(err.Error()))
	}

	discord.RespondToInteraction(s, i, fmt.Sprintf("[Video](%s) has been added successfully!", videoInfo.Link))
}
