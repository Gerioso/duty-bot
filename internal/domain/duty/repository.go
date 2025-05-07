package duty

import (
	"duty-bot/internal/domain/entities"
	"errors"
	"time"
)

var (
	ErrDutyNotFound = errors.New("дежурный не найден")
)

type DutyRepository interface {
	GetCurrentDuty() (entities.Duty, error)
	GetLastDuty() (entities.Duty, error)
	SetDuty(name string, weekStart time.Time) error
	RotateDuty() error
}
