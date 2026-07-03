package jobRegistry

import (
	"sync"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type StatusUpdate struct {
	JobId  uuid.UUID
	Status string
}

type Registry struct {
	mutex        sync.RWMutex
	jobs         map[uuid.UUID]*Job
	register     chan *Job
	unregister   chan uuid.UUID
	scheduler    *cron.Cron
	statusUpdate chan StatusUpdate
}

type Job struct {
	ID          uuid.UUID
	Title       string
	Description string
	Schedule    string
	Func        func()
	Status      string

	ScheduleId cron.EntryID
}
