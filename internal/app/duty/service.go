package duty

import "time"

// Сущность Duty
type Duty struct {
	Name      string
	WeekStart time.Time
}

// Интерфейс хранилища
type Storage interface {
	GetCurrentDuty() (Duty, error)
	SetDuty(name string, weekStart string) error
}

// Сервис дежурных
type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) GetCurrentDuty() (Duty, error) {
	return s.storage.GetCurrentDuty()
}

func (s *Service) SetDuty(name string, weekStart string) error {
	return s.storage.SetDuty(name, weekStart)
}
