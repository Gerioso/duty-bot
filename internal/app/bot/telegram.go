package bot

import (
	"duty-bot/internal/app/duty"
	"log"
	"time"

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
	u.Timeout = 60 // Telegram рекомендует 50-60 сек

	for {
		updates, err := b.api.GetUpdates(u)
		if err != nil {
			log.Println("Ошибка получения обновлений:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, update := range updates {
			if update.Message == nil || !update.Message.IsCommand() {
				continue
			}

			// Обновляем Offset, чтобы не обрабатывать одно и то же сообщение
			u.Offset = update.UpdateID + 1

			switch update.Message.Command() {
			case "duty":
				handleDutyCommand(b.api, update.Message, b.service)
			case "set_schedule":
				handleSetScheduleCommand(b.api, update, b.service)
			case "checklist":
				handleChecksCommand(b.api, update.Message)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
				b.api.Send(msg)
			}
		}
	}
}
