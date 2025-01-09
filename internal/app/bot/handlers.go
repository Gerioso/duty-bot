package bot

import (
	"duty-bot/internal/app/duty"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleDutyCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, service *duty.Service) {
	currentDuty, err := service.GetCurrentDuty()
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка: не удалось получить дежурного")
		bot.Send(msg)
		return
	}

	// Логируем текущий дежурный
	fmt.Printf("Текущий дежурный: %s (с %s)\n", currentDuty.Name, currentDuty.WeekStart.Format("02.01.2006"))

	response := fmt.Sprintf("Сейчас дежурит: %s (с %s)", currentDuty.Name, currentDuty.WeekStart.Format("02.01.2006"))
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	bot.Send(msg)
}

func handleSetScheduleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, service *duty.Service) {
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

	err := service.SetDuty(name, weekStart)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка сохранения расписания: "+err.Error()))
		return
	}

	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Расписание обновлено!"))
}

// handleChecksCommand обрабатывает команду /checks
func handleChecksCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Получаем текущую дату в нужном формате
	today := time.Now().Format("02.01.2006")

	// Формируем сообщение с чек-листом
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

	// Отправляем сообщение пользователю
	msg := tgbotapi.NewMessage(message.Chat.ID, checklist)
	bot.Send(msg)
}
