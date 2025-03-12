package duty_repo

import (
	"database/sql"
	"duty-bot/internal/domain/models"
	"time"
)

type SQLiteDutyRepository struct {
	db *sql.DB
}

func NewSQLiteDutyRepository(db *sql.DB) (*SQLiteDutyRepository, error) {
	_, err := db.Exec(`
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
	return &SQLiteDutyRepository{db: db}, nil
}

func (r *SQLiteDutyRepository) GetCurrentDuty() (models.Duty, error) {
	var d models.Duty
	row := r.db.QueryRow(`
        SELECT name, week_start FROM schedule 
        WHERE week_start <= date('now') 
        ORDER BY week_start DESC LIMIT 1;
    `)

	var weekStart string
	if err := row.Scan(&d.Name, &weekStart); err != nil {
		return d, err
	}

	d.WeekStart, _ = time.Parse("2006-01-02T15:04:05Z", weekStart)
	return d, nil
}
func (r *SQLiteDutyRepository) GetLastDuty() (models.Duty, error) {
	var d models.Duty
	row := r.db.QueryRow("SELECT name, week_start FROM schedule ORDER BY week_start DESC LIMIT 1;")

	var weekStart string
	if err := row.Scan(&d.Name, &weekStart); err != nil {
		return d, err
	}

	d.WeekStart, _ = time.Parse("2006-01-02T15:04:05Z", weekStart)
	return d, nil
}

func (r *SQLiteDutyRepository) SetDuty(name string, weekStart time.Time) error {
	_, err := r.db.Exec(`
        INSERT INTO schedule (name, week_start)
        VALUES (?, ?)
        ON CONFLICT(week_start) DO UPDATE SET name = excluded.name;
    `, name, weekStart)
	return err
}

func (r *SQLiteDutyRepository) RotateDuty() error {
	// Логика ротации пока не реализовал)
	return nil
}
