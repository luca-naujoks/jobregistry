package jobRegistry

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

var registry *Registry

func New() *Registry {
	r := &Registry{
		jobs:         make(map[uuid.UUID]*Job),
		register:     make(chan *Job),
		unregister:   make(chan uuid.UUID),
		scheduler:    cron.New(),
		statusUpdate: make(chan StatusUpdate),
	}

	go r.run()
	r.scheduler.Start()

	return r
}

func Init() {
	registry = New()
}

func (r *Registry) run() {
	for {
		select {
		case job := <-r.register:
			r.mutex.Lock()
			r.jobs[job.ID] = job

			parsedDuration, err := time.ParseDuration(job.Schedule)
			if err != nil {
				message := fmt.Sprintf("[Error] Parsing job Schedule: %s", err)
				fmt.Println(message)
				continue
			}
			duration := fmt.Sprintf("@every %s", parsedDuration)

			scheduleId, err := r.scheduler.AddFunc(duration, job.Func)
			if err != nil {
				message := fmt.Sprintf("[Error] Scheduling Job with id: %s error: %s", job.ID, err)
				fmt.Println(message)
				continue
			}
			job.ScheduleId = scheduleId
			r.mutex.Unlock()
		case jobId := <-r.unregister:
			r.mutex.Lock()
			job, ok := r.jobs[jobId]
			if ok {
				delete(r.jobs, jobId)
			}

			r.mutex.Unlock()

			if ok {
				r.scheduler.Remove(job.ScheduleId)
			}
		case status := <-r.statusUpdate:
			r.mutex.Lock()

			job, ok := r.jobs[status.JobId]
			if ok {
				job.Status = status.Status
			}

			r.mutex.Unlock()
		}
	}
}
