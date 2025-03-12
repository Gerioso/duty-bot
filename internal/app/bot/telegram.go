package bot

import (
	"duty-bot/internal/domain/duty"
	"duty-bot/internal/domain/employee"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	api             *tgbotapi.BotAPI
	dutyService     *duty.DutyService
	employeeService *employee.EmployeeService
}

func NewTelegramBot(token string, dutyService *duty.DutyService, employeeService *employee.EmployeeService) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	return &TelegramBot{
		api:             bot,
		dutyService:     dutyService,
		employeeService: employeeService,
	}
}

func (b *TelegramBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 // Таймаут для long polling

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "duty":
			handleDutyCommand(b.api, update.Message, b.dutyService)
		case "set_schedule":
			handleSetScheduleCommand(b.api, update, b.dutyService)
		// case "rotate":
		// 	handleRotateCommand(b.api, update.Message, b.dutyService)
		case "add_employee":
			handleAddEmployeeCommand(b.api, update.Message, b.employeeService)
		case "remove_employee":
			handleRemoveEmployeeCommand(b.api, update.Message, b.employeeService)
		case "checklist":
			handleChecksCommand(b.api, update.Message)
		case "help":
			handleHelpCommand(b.api, update.Message)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
			b.api.Send(msg)
		}
	}
}
