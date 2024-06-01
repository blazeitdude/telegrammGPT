package telegramBot

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"strings"
	"telegrammGPT/pkg/botLogger"
	"telegrammGPT/pkg/gptClient"
)

type BotConfiguration struct {
	ApiKey        string `yaml:"apiKey"`
	UpdateTimeout int    `yaml:"UpdateTimeout"`
	Retries       int    `yaml:"Retries"`
}

type TelegramBot struct {
	BotInstance  *tgbotapi.BotAPI
	UpdateConfig tgbotapi.UpdateConfig
	retries      int
}

func InitBot(conf BotConfiguration) TelegramBot {
	log := botLogger.GetLogger()
	bot, err := tgbotapi.NewBotAPI(conf.ApiKey)
	if err != nil {
		log.Logger.Fatalf("Failed to init Telegram Bot API Client", err)
	}
	log.Logger.Infof("Authorized on account %s", bot.Self.UserName)
	err = tgbotapi.SetLogger(log)
	if err != nil {
		log.Logger.Errorf("Failed to setUp logger to Telegram Bot API Client, error: ", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = conf.UpdateTimeout

	return TelegramBot{
		BotInstance:  bot,
		UpdateConfig: u,
		retries:      conf.Retries,
	}
}

func (b *TelegramBot) StartBot(client gptClient.GptClient) {
	log := botLogger.GetLogger()
	updates, err := b.BotInstance.GetUpdatesChan(b.UpdateConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to start receiving message from API", err)
	}
	log.Logger.Info("Bot started.")
	for update := range updates {
		if update.Message == nil {
			continue
		}
		text := update.Message.Text
		update.Message.Command()

		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID
		log.Logger.Debugf("message from [%s][%d] %s", update.Message.From.UserName, userID, text)

		if com := update.Message.Command(); com != "" {
			b.handleCommand(com, chatID)
			continue
		}
		gptResponse, err := client.SendMessage(text)
		retries := b.retries
		for gptResponse.ResponseBody == "" && retries > 0 {
			log.Logger.Debugf("retry to send request to ChatGPT [%d]", b.retries-retries)
			gptResponse, err = client.SendMessage(text)
			if err != nil {
				log.Logger.Debugf("Failed to send request ot ChatGPT: ", err)
			}
			retries--
		}
		msg := tgbotapi.NewMessage(chatID, gptResponse.ResponseBody)
		_, err = b.BotInstance.Send(msg)
		if err != nil {
			log.Logger.Debugf("Failed to send response to [%s][%d]", update.Message.From.UserName, userID)
		}
	}
}

func (b *TelegramBot) handleCommand(command string, chatID int64) bool {
	switch strings.ToLower(command) {
	case "/start":
		message := "Привет Некитка😏"
		tgbotapi.NewMessage(chatID, message)
		return true
	default:
		return false
	}
}
