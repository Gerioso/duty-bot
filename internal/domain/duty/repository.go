package duty

import (
	"duty-bot/internal/domain/models"
	"errors"
	"time"
)

var (
	ErrDutyNotFound = errors.New("дежурный не найден")
)

type DutyRepository interface {
	GetCurrentDuty() (models.Duty, error)
	SetDuty(name string, weekStart time.Time) error
	RotateDuty() error
}
