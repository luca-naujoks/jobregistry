package jobRegistry

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

var registry *Registry

func New() {
	r := &Registry{
		jobs:       make(map[uuid.UUID]*Job),
		register:   make(chan *Job),
		unregister: make(chan uuid.UUID),
		scheduler:  cron.New(),
		broadcast:  make(chan []byte),
	}

	go r.run()
	r.scheduler.Start()
}

func (r *Registry) run() {
	for {
		select {
		case job := <-r.register:
			r.muJobs.Lock()
			defer r.muJobs.Unlock()
			r.jobs[job.ID] = job

			parsedDuration, err := time.ParseDuration(job.Schedule)
			if err != nil {
				message := fmt.Sprintf("[Error] Parsing job Schedul: %s", err)
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
		case jobId := <-r.unregister:
			r.muJobs.Lock()
			defer r.muJobs.Unlock()
			job, ok := r.jobs[jobId]
			if !ok {
				message := fmt.Sprintf("[Warn] job with id: %s was not found inside the Registry", jobId)
				fmt.Println(message)
				continue
			}
			delete(r.jobs, jobId)
			r.scheduler.Remove(job.ScheduleId)
		}
	}
}

func GetById(id uuid.UUID) (Job, error) {
	registry.muJobs.RLocker()
	defer registry.muJobs.RUnlock()

	job, ok := registry.jobs[id]
	if !ok {
		errorMessage := fmt.Errorf("[Warn] job with id: %s was not found inside the Registry", id)
		return Job{}, errorMessage
	}
	return *job, nil
}

func GetAll() []Job {
	registry.muJobs.RLocker()
	defer registry.muJobs.RUnlock()

	jobs := make([]Job, len(registry.jobs))

	i := 0
	for _, job := range registry.jobs {
		jobs[i] = *job
		i++
	}

	return jobs
}
