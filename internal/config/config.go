package config

import (
	"log/slog"

	"github.com/joho/godotenv"

	"os"
)

type configuration struct {
	AppEnvironment      string
	BotPrefix           string
	BotStatus           string
	BotGuildJoinMessage string
	DiscordToken        string
	YoutubeApiKey       string
	RapidApiKey         string
}

const AppEnvironmentTest string = "TEST"

// config contains all environment variables that should be included in the .env file.
var config *configuration

func init() {
	Load()
}

// Load loads the environment variables from the .env file.
func Load() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("failed to load environment variables")
	}

	config = &configuration{
		AppEnvironment:      os.Getenv("APP_ENVIRONMENT"),
		BotPrefix:           os.Getenv("BOT_PREFIX"),
		BotStatus:           os.Getenv("BOT_STATUS"),
		BotGuildJoinMessage: os.Getenv("BOT_GUILD_JOIN_MESSAGE"),
		DiscordToken:        os.Getenv("DISCORD_TOKEN"),
		YoutubeApiKey:       os.Getenv("YOUTUBE_API_KEY"),
		RapidApiKey:         os.Getenv("RAPIDAPI_TIKTOK_API_KEY"),
	}
}

//func GetAppEnvironment() string {
//	return config.AppEnvironment
//}

func IsAppEnvironment(environments ...string) bool {
	if len(environments) == 0 {
		return config.AppEnvironment == environments[0]
	}

	for _, environment := range environments {
		if config.AppEnvironment == environment {
			return true
		}
	}

	return false
}

//func GetBotPrefix() string {
//	return config.BotPrefix
//}
//
//func GetBotStatus() string {
//	return config.BotStatus
//}
//
//func GetBotGuildJoinMessage() string {
//	return config.BotGuildJoinMessage
//}

func GetDiscordToken() string {
	return config.DiscordToken
}

func GetRapidApiKey() string {
	return config.RapidApiKey
}
