package config

import (
	"github.com/joho/godotenv"
	"log/slog"

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
	appEnv := os.Getenv("APP_ENVIRONMENT")
	if appEnv == "" || appEnv == "DEVELOPMENT" {
		// Load .env file only in development mode
		err := godotenv.Load()
		if err != nil {
			slog.Error("failed to load environment variables from .env file")
		}
	}

	config = &configuration{
		AppEnvironment: getEnvOrPanic("APP_ENVIRONMENT"),
		//BotPrefix:           getEnvOrPanic("BOT_PREFIX"),
		//BotStatus:           getEnvOrPanic("BOT_STATUS"),
		//BotGuildJoinMessage: getEnvOrPanic("BOT_GUILD_JOIN_MESSAGE"),
		DiscordToken: getEnvOrPanic("DISCORD_TOKEN"),
		RapidApiKey:  getEnvOrPanic("RAPIDAPI_API_KEY"),
	}
}

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		slog.Error("missing required environment variable", slog.String("key", key))
		panic("missing required environment variable: " + key)
	}
	return value
}

//	func GetAppEnvironment() string {
//		return config.AppEnvironment
//	}
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
