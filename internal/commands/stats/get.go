package stats

import (
	"clipify-bot/internal/discord"
	"clipify-bot/internal/firebase"
	"clipify-bot/internal/models"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sort"
	"strconv"
	"strings"
)

func GetStats(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	username := i.Member.User.Username
	accounts, err := firebase.GetUserAccounts(ctx, username)
	if err != nil || accounts == nil {
		discord.RespondToInteractionEmbedError(s, i, fmt.Sprintf("No accounts found for **%s**", username))
		return
	}

	page := 0

	videos, _ := firebase.GetAllAccountVideos(ctx, username, accounts[page].Username, accounts[page].Platform)

	embed := buildAccountStatsEmbed(accounts[page].Username, accounts[page].Platform, videos)
	components := buildAccountNavComponents(page, len(accounts))

	discord.RespondToInteractionEmbedAndButtons(s, i, embed, components)
}

func HandleAccountStatsButton(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	username := i.Member.User.Username
	accounts, err := firebase.GetUserAccounts(ctx, username)
	if err != nil || accounts == nil {
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
	} else if action == "next" && page < len(accounts)-1 {
		page++
	}

	videos, _ := firebase.GetAllAccountVideos(ctx, username, accounts[page].Username, accounts[page].Platform)

	embed := buildAccountStatsEmbed(accounts[page].Username, accounts[page].Platform, videos)
	components := buildAccountNavComponents(page, len(accounts))

	discord.RespondToButtonInteractionEmbedAndButtons(s, i, embed, components)
}

func buildAccountStatsEmbed(accountName string, platform models.Platform, videos []models.Video) *discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField
	var totalViews = 0

	fields = append(fields, &discordgo.MessageEmbedField{
		Inline: false,
	})
	if len(videos) == 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s - %s", platform, accountName),
			Value:  "No videos found.",
			Inline: false,
		})
	} else {
		fields = append(fields, &discordgo.MessageEmbedField{
			Inline: false,
		})

		sort.Slice(videos, func(i, j int) bool {
			return videos[i].Views > videos[j].Views
		})
		for idx, video := range videos {
			totalViews += video.Views
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
	fields[0].Value = fmt.Sprintf("Total views - %s", formatKM(totalViews))

	return &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("**📊 Stats**\n%s (%s)", accountName, platform),
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

func formatKM(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}
