package duty

import (
	"duty-bot/internal/domain/models"
	"errors"
	"time"
)

type DutyService struct {
	repo DutyRepository
}

func NewDutyService(repo DutyRepository) *DutyService {
	return &DutyService{repo: repo}
}

func (s *DutyService) GetCurrentDuty() (models.Duty, error) {
	return s.repo.GetCurrentDuty()
}

func (s *DutyService) SetDuty(name string, weekStart time.Time) error {
	if name == "" {
		return errors.New("имя не может быть пустым")
	}
	return s.repo.SetDuty(name, weekStart)
}

func (s *DutyService) RotateDuty() error {
	return s.repo.RotateDuty()
}
