package bot

import (
	"duty-bot/internal/domain/duty"
	"duty-bot/internal/domain/employee"
	"duty-bot/internal/domain/entities"
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
		case entities.CommandDuty:
			handleDutyCommand(b.api, update.Message, b.dutyService)
		case entities.CommandSetSchedule:
			handleSetScheduleCommand(b.api, update, b.dutyService)
		case entities.CommandAddEmployee:
			handleAddEmployeeCommand(b.api, update.Message, b.employeeService)
		case entities.CommandRemoveEmployee:
			handleRemoveEmployeeCommand(b.api, update.Message, b.employeeService)
		case entities.CommandChecklist:
			handleChecksCommand(b.api, update.Message)
		case entities.CommandHelp:
			handleHelpCommand(b.api, update.Message)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
			b.api.Send(msg)
		}
	}
}
