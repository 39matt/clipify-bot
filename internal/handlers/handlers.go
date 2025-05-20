package handlers

import (
	"clipping-bot/internal/commands/account"
	"clipping-bot/internal/commands/verification"
	"context"
	"github.com/bwmarrin/discordgo"
	"time"
)

func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch i.ApplicationCommandData().Name {
	case "register":
		account.Register(ctx, s, i)
	case "add-account":
		account.AddAccount(ctx, s, i)
	case "verify-account":
		verification.VerifyAccount(ctx, s, i)
	case "remove-verification":
		verification.RemoveVerification(ctx, s, i)
	}
}
