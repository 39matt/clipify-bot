package stats

import (
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"clipping-bot/internal/models"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sort"
	"strconv"
	"strings"
)

func GetStats(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	username := i.Member.User.Username
	accountNames, err := firebase.GetUserAccountNames(ctx, username)
	if err != nil || accountNames == nil {
		discord.RespondToInteractionEmbedError(s, i, fmt.Sprintf("No accounts found for **%s**", username))
		return
	}

	page := 0
	platform := "TikTok"

	videos, _ := firebase.GetAllAccountVideos(ctx, username, accountNames[page], platform)

	embed := buildAccountStatsEmbed(username, accountNames[page], platform, videos)
	components := buildAccountNavComponents(page, len(accountNames))

	discord.RespondToInteractionEmbedAndButtons(s, i, embed, components)
}

func HandleAccountStatsButton(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	username := i.Member.User.Username
	accountNames, err := firebase.GetUserAccountNames(ctx, username)
	if err != nil || accountNames == nil {
		discord.RespondToInteractionEmbedError(s, i, fmt.Sprintf("No accounts found for **%s**", username))
		return
	}

	data := i.MessageComponentData()
	var action string
	var page int
	parts := strings.Split(data.CustomID, "_")
	if len(parts) != 4 {
		discord.RespondToInteractionEmbedError(s, i, "Invalid button.")
		return
	}
	action = parts[2]
	page, _ = strconv.Atoi(parts[3])

	if action == "prev" && page > 0 {
		page--
	} else if action == "next" && page < len(accountNames)-1 {
		page++
	}

	platform := "TikTok"
	videos, _ := firebase.GetAllAccountVideos(ctx, username, accountNames[page], platform)

	embed := buildAccountStatsEmbed(username, accountNames[page], platform, videos)
	components := buildAccountNavComponents(page, len(accountNames))

	discord.RespondToInteractionEmbedAndButtons(s, i, embed, components)
}

func buildAccountStatsEmbed(username, accountName, platform string, videos []models.Video) *discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	if len(videos) == 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s - %s", platform, accountName),
			Value:  "No videos found.",
			Inline: false,
		})
	} else {
		sort.Slice(videos, func(i, j int) bool {
			return videos[i].Views > videos[j].Views
		})
		for idx, video := range videos {
			name := video.Name
			if len(name) > 30 {
				name = name[:27] + "..."
			}
			value := fmt.Sprintf(
				"\n[🔗 %s](<%s>)\n👁️ **%s** • ❤️ **%s** • 💬 **%s** • 🔄 **%s**\n",
				video.Name,
				video.Link,
				formatNumber(video.Views),
				formatNumber(video.Likes),
				formatNumber(video.Comments),
				formatNumber(video.Shares),
			)
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("----------VIDEO %d-----------", idx+1),
				Value:  value,
				Inline: false,
			})
		}
	}

	return &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("**📊 Stats**"),
		Color:  0x50C878,
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Use /add-video to add more videos!",
		},
	}
}

func buildAccountNavComponents(page, total int) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Previous",
					Style:    discordgo.SecondaryButton,
					CustomID: fmt.Sprintf("account_stats_prev_%d", page),
					Disabled: page == 0,
				},
				discordgo.Button{
					Label:    "Next",
					Style:    discordgo.SecondaryButton,
					CustomID: fmt.Sprintf("account_stats_next_%d", page),
					Disabled: page == total-1,
				},
			},
		},
	}
}

func formatNumber(n int) string {
	str := fmt.Sprintf("%d", n)
	if n < 1000 {
		return str
	}
	result := ""
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(c)
	}
	return result
}
