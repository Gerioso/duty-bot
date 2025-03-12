package duty

import (
	"database/sql"
	"duty-bot/internal/domain/employee"
	"duty-bot/internal/domain/models"
	"errors"
	"time"
)

type DutyService struct {
	dutyRepo     DutyRepository
	employeeRepo employee.EmployeeRepository
}

func NewDutyService(dutyRepo DutyRepository, employeeRepo employee.EmployeeRepository) *DutyService {
	return &DutyService{
		dutyRepo:     dutyRepo,
		employeeRepo: employeeRepo,
	}
}

func (s *DutyService) GetCurrentDuty() (models.Duty, error) {
	d, err := s.dutyRepo.GetCurrentDuty()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := s.RotateDuty(); err != nil {
				return d, err
			}
			return s.dutyRepo.GetCurrentDuty()
		}
		return d, err
	}

	return d, nil
}

func (s *DutyService) SetDuty(name string, weekStart time.Time) error {
	if name == "" {
		return errors.New("имя не может быть пустым")
	}
	return s.dutyRepo.SetDuty(name, weekStart)
}

func (s *DutyService) RotateDuty() error {

	employees, err := s.employeeRepo.GetEmployees()
	if err != nil {
		return err
	}

	if len(employees) == 0 {
		return errors.New("нет сотрудников для ротации")
	}

	lastDuty, err := s.dutyRepo.GetLastDuty()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	lastIndex := -1
	for i, name := range employees {
		if name == lastDuty.Name {
			lastIndex = i
			break
		}
	}

	currentIndex := (lastIndex + 1) % len(employees)
	currentDuty := employees[currentIndex]

	nextIndex := (currentIndex + 1) % len(employees)
	nextDuty := employees[nextIndex]

	now := time.Now()
	currentWeekStart := now.AddDate(0, 0, -int(now.Weekday()))
	if err := s.dutyRepo.SetDuty(currentDuty, currentWeekStart); err != nil {
		return err
	}

	// Назначаем дежурного на следующую неделю
	nextWeekStart := now.AddDate(0, 0, 7-int(now.Weekday()))
	if err := s.dutyRepo.SetDuty(nextDuty, nextWeekStart); err != nil {
		return err
	}

	return nil
}
