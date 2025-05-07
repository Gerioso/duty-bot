package entities

import "time"

type Duty struct {
	Name      string
	WeekStart time.Time
}

const (
	CommandDuty           = "duty"
	CommandSetSchedule    = "setSchedule"
	CommandAddEmployee    = "addEmployee"
	CommandRemoveEmployee = "removeEmployee"
	CommandChecklist      = "checklist"
	CommandHelp           = "help"
)
