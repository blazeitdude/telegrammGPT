package main

import (
	"telegrammGPT/cmd/config"
	"telegrammGPT/pkg/botLogger"
	gptClient2 "telegrammGPT/pkg/gptClient"
	tgbot "telegrammGPT/pkg/telegramBot"
)

func main() {
	conf := config.ReadConfig()
	botLogger.InitLogger(conf.LoggerConfig)
	telegramBot := tgbot.InitBot(conf.BotConfig)
	gptClient := gptClient2.InitGpt(conf.GptConfig)
	telegramBot.StartBot(gptClient)
}
