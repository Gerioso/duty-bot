package bot

import (
	"duty-bot/internal/domain/duty"
	"duty-bot/internal/domain/employee"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleDutyCommand обрабатывает команду /duty
func handleDutyCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dutyService *duty.DutyService) {
	currentDuty, err := dutyService.GetCurrentDuty()
	if err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Ошибка: "+err.Error()))
		return
	}

	response := fmt.Sprintf("Сейчас дежурит: %s (с %s)", currentDuty.Name, currentDuty.WeekStart.Format("02.01.2006"))
	bot.Send(tgbotapi.NewMessage(message.Chat.ID, response))
}

// handleSetScheduleCommand обрабатывает команду /set_schedule
func handleSetScheduleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, dutyService *duty.DutyService) {
	if update.Message == nil || update.Message.Chat.Type != "private" {
		return
	}

	args := strings.Fields(update.Message.Text)
	if len(args) != 3 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Формат: /set_schedule <имя> <дата начала недели в формате YYYY-MM-DD>"))
		return
	}

	name := args[1]
	weekStart := args[2]

	parsedWeekStart, err := time.Parse("2006-01-02", weekStart)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка: неверный формат даты. Используйте YYYY-MM-DD"))
		return
	}

	if err := dutyService.SetDuty(name, parsedWeekStart); err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка: "+err.Error()))
		return
	}

	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Расписание обновлено!"))
}

// handleRotateCommand обрабатывает команду /rotate
func handleRotateCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dutyService *duty.DutyService) {
	if err := dutyService.RotateDuty(); err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Ошибка: "+err.Error()))
		return
	}

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Дежурные успешно ротированы!"))
}

// handleAddEmployeeCommand обрабатывает команду /add_employee
func handleAddEmployeeCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, employeeService *employee.EmployeeService) {
	args := strings.Fields(message.Text)
	if len(args) != 2 {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Формат: /add_employee <имя>"))
		return
	}

	name := args[1]
	if err := employeeService.AddEmployee(name); err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Ошибка: "+err.Error()))
		return
	}

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Сотрудник добавлен!"))
}

// handleRemoveEmployeeCommand обрабатывает команду /remove_employee
func handleRemoveEmployeeCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, employeeService *employee.EmployeeService) {
	args := strings.Fields(message.Text)
	if len(args) != 2 {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Формат: /remove_employee <имя>"))
		return
	}

	name := args[1]
	if err := employeeService.RemoveEmployee(name); err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Ошибка: "+err.Error()))
		return
	}

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Сотрудник удалён!"))
}

// handleChecksCommand обрабатывает команду /checks
func handleChecksCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	today := time.Now().Format("02.01.2006")

	checklist := fmt.Sprintf(`Шаблон Чек-листа для дежурного на %s

Состояние стенда 💼
Оплата💸
 1) Оплата новой картой  ✅
 2) Оплата с привязкой карты  ✅
 3) Оплата картой с 3DS v2 ❌ - MNT-1217
 4) Оплата СБП ✅

Привязка карт 📌
 1) Создание привязки карты ✅
 2) Удаление привязки карты ✅

Капча 🧩
 1) Автокапча ✅
 2) Ручная капча ✅

Управление платежами 🔄
 1) Отмена платежа ✅
 2) Реверс платежа ✅
 3) Возврат платежа ✅
`, today)

	msg := tgbotapi.NewMessage(message.Chat.ID, checklist)
	bot.Send(msg)
}

// handleHelpCommand обрабатывает команду /help
func handleHelpCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	helpText := `Доступные команды:
/duty - Показать текущего дежурного
/set_schedule <имя> <дата> - Установить расписание дежурств
/rotate - Ротировать дежурных
/add_employee <имя> - Добавить сотрудника
/remove_employee <имя> - Удалить сотрудника
/checks - Показать чек-лист
/help - Показать это сообщение`

	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	bot.Send(msg)
}
