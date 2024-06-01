package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"telegrammGPT/pkg/botLogger"
	"telegrammGPT/pkg/gptClient"
	"telegrammGPT/pkg/telegramBot"
)

type Configuration struct {
	LoggerConfig botLogger.LoggerConfig       `yaml:"loggerConfig"`
	BotConfig    telegramBot.BotConfiguration `yaml:"botConfig"`
	GptConfig    gptClient.GptConfiguration   `yaml:"gptConfig"`
}

func ReadConfig() Configuration {
	var conf Configuration
	viper.SetConfigName("startUp")
	viper.SetConfigType("yml")
	viper.AddConfigPath("../../config")
	if err := viper.ReadInConfig(); err != nil {
		var configParseError viper.ConfigParseError
		var configFileNotFoundError viper.ConfigFileNotFoundError
		var configMarshalError viper.ConfigMarshalError
		var unsupportedConfigError viper.UnsupportedConfigError
		switch {
		case errors.As(err, &configParseError):
			log.Fatal("Failed to parse config file")
		case errors.As(err, &configFileNotFoundError):
			log.Fatal("Config file not found")
		case errors.As(err, &configMarshalError):
			log.Fatal("Failed to marshall config file")
		case errors.As(err, &unsupportedConfigError):
			log.Fatal("This config file is unsupported")
		default:
			log.Fatal("unexpected error while reading the configuration file")
		}
	}
	conf.BotConfig.ApiKey = os.Getenv("TELEGRAMM_APIKEY")
	conf.GptConfig.ApiKey = os.Getenv("GPT_APIKEY")

	if conf.BotConfig.ApiKey == "" || conf.GptConfig.ApiKey == "" {
		log.Fatal("Fail;ed to read ENV")
	}

	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal("Failed to unmarshall config file")

	}
	return conf
}
