package jobRegistry

import (
	"sync"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type Entry struct {
	Job    Job
	CronId uuid.UUID
}
type Registry struct {
	muJobs     sync.RWMutex
	jobs       map[uuid.UUID]*Job
	register   chan *Job
	unregister chan uuid.UUID
	scheduler  *cron.Cron
	broadcast  chan []byte
}

type Job struct {
	ID          uuid.UUID
	Title       string
	Description string
	Schedule    string
	Func        func()

	ScheduleId cron.EntryID
}
