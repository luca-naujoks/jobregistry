package jobRegistry

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func GetById(id uuid.UUID) (Job, error) {
	registry.mutex.RLock()
	defer registry.mutex.RUnlock()

	job, ok := registry.jobs[id]
	if !ok {
		errorMessage := fmt.Errorf("[Warn] job with id: %s was not found inside the Registry", id)
		return Job{}, errorMessage
	}
	return *job, nil
}

func Get() []Job {
	registry.mutex.RLock()
	defer registry.mutex.RUnlock()

	jobs := make([]Job, len(registry.jobs))

	i := 0
	for _, job := range registry.jobs {
		jobs[i] = *job
		i++
	}

	return jobs
}

func Register(job Job) error {
	registry.register <- &job
	return nil
}

func Remove(id uuid.UUID) error {
	registry.unregister <- id
	return nil
}

func Execute(id uuid.UUID) error {
	job, ok := registry.jobs[id]
	if !ok {
		message := fmt.Sprintf("[Error] Job with id: %s was not found inside the registry", id.String())
		return errors.New(message)
	}
	registry.scheduler.Entry(job.ScheduleId).Job.Run()

	return nil
}
