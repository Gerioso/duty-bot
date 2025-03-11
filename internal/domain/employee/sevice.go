package employee

import "errors"

type EmployeeService struct {
	repo EmployeeRepository
}

func NewEmployeeService(repo EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

func (s *EmployeeService) AddEmployee(name string) error {
	if name == "" {
		return errors.New("имя не может быть пустым")
	}
	return s.repo.AddEmployee(name)
}

func (s *EmployeeService) RemoveEmployee(name string) error {
	if name == "" {
		return errors.New("имя не может быть пустым")
	}
	return s.repo.RemoveEmployee(name)
}

func (s *EmployeeService) GetEmployees() ([]string, error) {
	return s.repo.GetEmployees()
}
