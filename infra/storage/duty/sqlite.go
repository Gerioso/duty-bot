package duty_repo

import (
	"database/sql"
	"duty-bot/internal/domain/entities"
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

func (r *SQLiteDutyRepository) GetCurrentDuty() (entities.Duty, error) {
	var d entities.Duty
	now := time.Now()
	currentWeekStart := now.AddDate(0, 0, -int(now.Weekday())).Format("2006-01-02")
	row := r.db.QueryRow(`
        SELECT name, week_start FROM schedule 
        WHERE DATE(week_start) = ?;
    `, currentWeekStart)

	var weekStart string
	if err := row.Scan(&d.Name, &weekStart); err != nil {
		return d, err
	}

	d.WeekStart, _ = time.Parse("2006-01-02T15:04:05Z", weekStart)
	return d, nil
}
func (r *SQLiteDutyRepository) GetLastDuty() (entities.Duty, error) {
	var d entities.Duty
	row := r.db.QueryRow("SELECT name, week_start FROM schedule ORDER BY week_start DESC LIMIT 1;")

	var weekStart string
	if err := row.Scan(&d.Name, &weekStart); err != nil {
		return d, err
	}

	d.WeekStart, _ = time.Parse("2006-01-02T15:04:05Z", weekStart)
	return d, nil
}

func (r *SQLiteDutyRepository) SetDuty(name string, weekStart time.Time) error {
	weekStartStr := weekStart.Format("2006-01-02")

	_, err := r.db.Exec(`
        INSERT INTO schedule (name, week_start)
        VALUES (?, ?)
        ON CONFLICT(week_start) DO UPDATE SET name = excluded.name;
    `, name, weekStartStr)
	return err
}

func (r *SQLiteDutyRepository) RotateDuty() error {
	// Логика ротации пока не реализовал)
	return nil
}
