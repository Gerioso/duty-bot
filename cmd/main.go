package main

import (
	"database/sql"
	dutyRepo "duty-bot/infra/storage/duty"
	employeeRepo "duty-bot/infra/storage/employee"
	"duty-bot/internal/app/bot"
	"duty-bot/internal/app/config"
	"duty-bot/internal/domain/duty"
	"duty-bot/internal/domain/employee"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	db, err := sql.Open("sqlite3", cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dutyRepo, err := dutyRepo.NewSQLiteDutyRepository(db)
	if err != nil {
		log.Fatalf("Ошибка cоздания бд дежурств: %v", err)
	}
	employeeRepo, err := employeeRepo.NewSQLiteEmployeeRepository(db)
	if err != nil {
		log.Fatalf("Ошибка cоздания бд сотрудников: %v", err)
	}
	dutyService := duty.NewDutyService(dutyRepo)
	employeeService := employee.NewEmployeeService(employeeRepo)

	// Инициализируем Telegram-бота
	telegramBot := bot.NewTelegramBot(cfg.TelegramToken, dutyService, employeeService)

	// Запускаем бота
	log.Println("Запуск Telegram-бота...")
	telegramBot.Start()
}
