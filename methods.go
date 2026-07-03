package registry

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *Registry) Job(id uuid.UUID) (Job, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	job, ok := r.jobs[id]
	if !ok {
		errorMessage := fmt.Errorf("[Warn] job with id: %s was not found inside the Registry", id)
		return Job{}, errorMessage
	}
	return *job, nil
}

func (r *Registry) Jobs() []Job {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	jobs := make([]Job, 0, len(r.jobs))

	for _, job := range r.jobs {
		jobs = append(jobs, *job)
	}

	return jobs
}

func (r *Registry) Register(job Job) error {
	parsedDuration, err := time.ParseDuration(job.Schedule)
	if err != nil {
		return fmt.Errorf("[Error] Parsing job Schedule: %w", err)
	}
	duration := fmt.Sprintf("@every %s", parsedDuration)

	scheduleId, err := r.scheduler.AddFunc(duration, job.Func)
	if err != nil {
		return fmt.Errorf("[Error] Scheduling Job with id: %s error: %w", job.ID, err)
	}
	job.ScheduleId = scheduleId

	r.mutex.Lock()
	r.jobs[job.ID] = &job
	r.mutex.Unlock()

	return nil
}

func (r *Registry) Remove(id uuid.UUID) error {
	r.mutex.Lock()
	job, ok := r.jobs[id]
	if ok {
		delete(r.jobs, id)
	}
	r.mutex.Unlock()

	if !ok {
		return fmt.Errorf("job %s not found", id)
	}

	r.scheduler.Remove(job.ScheduleId)

	return nil
}

func (r *Registry) Run(id uuid.UUID) error {
	r.mutex.RLock()
	job, ok := r.jobs[id]
	r.mutex.RUnlock()
	if !ok {
		message := fmt.Sprintf("[Error] Job with id: %s was not found inside the registry", id.String())
		return errors.New(message)
	}
	r.scheduler.Entry(job.ScheduleId).Job.Run()

	return nil
}

func (r *Registry) statusUpdate(status StatusUpdate) error {
	return nil
}
