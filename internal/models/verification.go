package models

import "time"

type PendingVerification struct {
	Code        int       `json:"code"`
	Platform    string    `json:"platform"`
	AccountName string    `json:"accountname"`
	DiscordID   string    `json:"discordid"`
	CreatedAt   time.Time `json:"createdat"`
}
