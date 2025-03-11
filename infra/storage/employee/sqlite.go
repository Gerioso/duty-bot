package employee

import (
	"database/sql"
)

type SQLiteEmployeeRepository struct {
	db *sql.DB
}

func NewSQLiteEmployeeRepository(db *sql.DB) (*SQLiteEmployeeRepository, error) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS employees (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE
	);
	`)
	if err != nil {
		return nil, err
	}
	return &SQLiteEmployeeRepository{db: db}, nil
}

func (r *SQLiteEmployeeRepository) AddEmployee(name string) error {
	_, err := r.db.Exec("INSERT INTO employees (name) VALUES (?);", name)
	return err
}

func (r *SQLiteEmployeeRepository) RemoveEmployee(name string) error {
	_, err := r.db.Exec("DELETE FROM employees WHERE name = ?;", name)
	return err
}

func (r *SQLiteEmployeeRepository) GetEmployees() ([]string, error) {
	rows, err := r.db.Query("SELECT name FROM employees;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		employees = append(employees, name)
	}

	return employees, nil
}
