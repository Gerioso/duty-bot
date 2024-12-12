package bot

import (
	"duty-bot/internal/app/duty"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	api     *tgbotapi.BotAPI
	service *duty.Service
}

func NewTelegramBot(token string, service *duty.Service) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	return &TelegramBot{api: bot, service: service}
}

func (b *TelegramBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Command() {
			case "duty":
				handleDutyCommand(b.api, update.Message, b.service)
			case "set_schedule":
				handleSetScheduleCommand(b.api, update, b.service)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
				b.api.Send(msg)
			}
		}
	}
}
