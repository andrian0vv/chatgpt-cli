package config

import "os"

type Config struct {
	OpenaiApiKey string
}

func New() Config {
	return Config{
		OpenaiApiKey: os.Getenv("OPENAI_API_KEY"),
	}
}
