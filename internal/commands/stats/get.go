package stats

import (
	"clipping-bot/internal/discord"
	"clipping-bot/internal/firebase"
	"clipping-bot/internal/models"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"strconv"
	"strings"
)

func GetStats(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	username := i.Member.User.Username
	accountNames, err := firebase.GetUserAccountNames(ctx, username)
	if err != nil || len(accountNames) == 0 {
		discord.RespondToInteraction(s, i, "No accounts found.")
		return
	}

	page := 0            // Start at first account
	platform := "TikTok" // Adjust if you support multiple platforms

	videos, _ := firebase.GetAllAccountVideos(ctx, username, accountNames[page], platform)
	embed := buildAccountStatsEmbed(username, accountNames[page], platform, videos)

	components := buildAccountNavComponents(page, len(accountNames))

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &components,
	})
	if err != nil {
		slog.Error("failed to respond to interaction: ", err)
		return
	}
}

func HandleAccountStatsButton(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	username := i.Member.User.Username
	accountNames, err := firebase.GetUserAccountNames(ctx, username)
	if err != nil || len(accountNames) == 0 {
		discord.RespondToInteraction(s, i, "No accounts found.")
		return
	}

	data := i.MessageComponentData()
	var action string
	var page int
	// CustomID format: account_stats_prev_0 or account_stats_next_1
	parts := strings.Split(data.CustomID, "_")
	if len(parts) != 4 {
		discord.RespondToInteraction(s, i, "Invalid button.")
		return
	}
	action = parts[2]
	page, _ = strconv.Atoi(parts[3])

	// Determine new page
	if action == "prev" && page > 0 {
		page--
	} else if action == "next" && page < len(accountNames)-1 {
		page++
	}

	platform := "TikTok" // Adjust if you support multiple platforms
	videos, _ := firebase.GetAllAccountVideos(ctx, username, accountNames[page], platform)
	embed := buildAccountStatsEmbed(username, accountNames[page], platform, videos)
	components := buildAccountNavComponents(page, len(accountNames))

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})

	if err != nil {
		slog.Error("failed to respond", "error", err)
	}
}

func buildAccountStatsEmbed(username, accountName, platform string, videos []models.Video) *discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{}

	if len(videos) == 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s - %s", platform, accountName),
			Value:  "No videos found.",
			Inline: false,
		})
	} else {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s - %s's Stats", platform, accountName),
			Value:  "",
			Inline: false,
		})
		for idx, video := range videos {
			name := video.Name
			if len(name) > 30 {
				name = name[:27] + "..."
			}
			value := fmt.Sprintf(
				"[🔗 Link](<%s>)\n👁️ **%s** • ❤️ **%s** • 💬 **%s** • 🔄 **%s**",
				video.Link,
				formatNumber(video.Views),
				formatNumber(video.Likes),
				formatNumber(video.Comments),
				formatNumber(video.Shares),
			)
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%d. %s", idx+1, name),
				Value:  value,
				Inline: false,
			})
		}
	}

	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s's Account: %s (%s)", username, accountName, platform),
		Description: "Here are your video stats for this account:",
		Color:       0x5865F2,
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Use /add-video to add more videos!",
		},
	}
}

// Helper: Build navigation buttons
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

// Helper: Format numbers with commas
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
