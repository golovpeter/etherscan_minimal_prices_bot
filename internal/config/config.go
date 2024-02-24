package config

import "os"

type Config struct {
	BotToken   string
	WebHookURL string
	ApiToken   string
	Port       string
}

func NewConfig() *Config {
	return &Config{
		BotToken:   os.Getenv("BOT_TOKEN"),
		WebHookURL: os.Getenv("WEBHOOK_URL"),
		ApiToken:   os.Getenv("API_KEY"),
		Port:       os.Getenv("PORT"),
	}
}
