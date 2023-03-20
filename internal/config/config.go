package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	TgBotToken string
	ApiKey     string
	Post       PostgresConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func MustLoad() Config {
	tgBotTokenToken := os.Getenv("BOT_T")
	apiKey := os.Getenv("API_T")
	if tgBotTokenToken == "" {
		log.Fatal("Bot token is not specified")
	}
	if apiKey == "" {
		log.Fatal("ApiKey token is not specified")
	}

	return Config{
		TgBotToken: tgBotTokenToken,
		ApiKey:     apiKey,
		Post: PostgresConfig{
			Host:     viper.GetString("postgres.host"),
			Port:     viper.GetString("postgres.port"),
			Username: viper.GetString("postgres.username"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   viper.GetString("postgres.dbname"),
			SSLMode:  viper.GetString("postgres.sslmode"),
		},
	}
}
