package employee

import "errors"

var (
	ErrEmployeeNotFound = errors.New("сотрудник не найден")
)

type EmployeeRepository interface {
	AddEmployee(name string) error
	RemoveEmployee(name string) error
	GetEmployees() ([]string, error)
}
