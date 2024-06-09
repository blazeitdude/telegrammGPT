package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
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
	viper.AddConfigPath(filepath.Join("config"))
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigParseError:
			log.Fatal("Failed to parse config file")
		case viper.ConfigFileNotFoundError:
			log.Fatal("Config file not found")
		case viper.ConfigMarshalError:
			log.Fatal("Failed to marshall config file")
		case viper.UnsupportedConfigError:
			log.Fatal("This config file is unsupported")
		default:
			log.Fatal("unexpected error while reading the configuration file")
		}
	}
	conf.BotConfig.ApiKey = os.Getenv("TELEGRAMM_APIKEY")
	conf.GptConfig.ApiKey = os.Getenv("GPT_APIKEY")

	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal("Failed to unmarshall config file")

	}
	return conf
}
