package storage

import (
	"database/sql"
	"duty-bot/internal/app/duty"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"time"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(path string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS schedule (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		week_start DATE,
		UNIQUE(week_start)
	);
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) GetCurrentDuty() (duty.Duty, error) {
	var d duty.Duty
	row := s.db.QueryRow(`
	SELECT name, week_start FROM schedule 
	WHERE week_start <= date('now') 
	ORDER BY week_start DESC LIMIT 1;
	`)

	var weekStart string
	err := row.Scan(&d.Name, &weekStart)
	if err != nil {
		// Логируем ошибку для диагностики
		fmt.Printf("Ошибка получения текущего дежурного: %v\n", err)
		return d, err
	}

	// Логируем дату, которую мы извлекли из базы данных
	fmt.Printf("Извлеченная дата начала недели: %s\n", weekStart)

	// Парсим дату из строки
	d.WeekStart, err = time.Parse("2006-01-02T15:04:05Z", weekStart)
	if err != nil {
		// Логируем ошибку парсинга даты
		fmt.Printf("Ошибка парсинга даты: %v\n", err)
		return d, err
	}

	return d, nil
}

func (s *SQLiteStorage) SetDuty(name string, weekStart string) error {
	// Логируем данные перед вставкой
	fmt.Printf("Запись дежурного: %s, дата начала недели: %s\n", name, weekStart)

	_, err := s.db.Exec(`
		INSERT INTO schedule (name, week_start)
		VALUES (?, ?)
		ON CONFLICT(week_start) DO UPDATE SET name = excluded.name;
	`, name, weekStart)
	return err
}
