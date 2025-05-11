package handlers

import (
	"clipping-bot/internal/commands/account"
	"clipping-bot/internal/config"
	"clipping-bot/internal/discord"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"strings"
)

// GuildCreateHandler will be called every time a new guild is joined.
func GuildCreateHandler(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, err := s.ChannelMessageSend(channel.ID, config.GetBotGuildJoinMessage())
			if err != nil {
				slog.Warn("failed to send guild create handler message", "error", err)
			}

			return
		}
	}
}

// MessageCreateHandler will be called everytime a new message is sent in a channel the bot has access to.
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { // Preventing bot from using own commands
		return
	}

	slog.Info("processing command", "command", m.Content)

	prefix := config.GetBotPrefix()
	//guildID := discord.SearchGuildByChannelID(m.ChannelID)
	//v := music.VoiceInstances[guildID]
	cmd := strings.Split(m.Content, " ") //	Splitting command into string slice
	fmt.Print(cmd)

	switch cmd[0] {
	case prefix + "help":
		discord.SendChannelMessage(m.ChannelID, "pomagaj!!")
	case prefix + "ping":
		discord.SendChannelMessage(m.ChannelID, "Pong!")
	//case prefix + "add-account":
	//	account.AddAccount(m.ChannelID, cmd[1], cmd[2])
	//case prefix + "dice":
	//	commands.RollDice(cmd, m)
	//case prefix + "insult":
	//	commands.PostInsult(m)
	//case prefix + "advice":
	//	commands.PostAdvice(m)
	//case prefix + "kanye":
	//	commands.PostKanyeQuote(m)
	//case prefix + "play":
	//	music.PlayMusic(cmd[1:], v, s, m)
	//case prefix + "leave":
	//	music.LeaveVoice(v, m)
	//case prefix + "skip":
	//	music.SkipMusic(v, m)
	//case prefix + "stop":
	//	music.StopMusic(v, m)
	//case prefix + "chess":
	//	chess.Menu(cmd[1:], s, m)
	default:
		return
	}
}

// Handles slash commands
func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "add-account":
		account.AddAccount(s, i)
	}
}
