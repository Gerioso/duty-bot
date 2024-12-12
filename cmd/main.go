package main

import (
	"duty-bot/internal/app/bot"
	"duty-bot/internal/app/config"
	"duty-bot/internal/app/duty"
	"duty-bot/internal/app/storage"
	"log"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализируем хранилище
	db, err := storage.NewSQLiteStorage(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Создаем сервис для работы с дежурными
	dutyService := duty.NewService(db)

	// Инициализируем Telegram-бота
	telegramBot := bot.NewTelegramBot(cfg.TelegramToken, dutyService)

	// Запускаем бота
	log.Println("Запуск Telegram-бота...")
	telegramBot.Start()
}
