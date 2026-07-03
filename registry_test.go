package jobRegistry

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRegisterJob(t *testing.T) {
	// Create and Start a New Registry
	r := New()

	// Set a deadline for the Job to register later cause it's asynchronous
	deadline := time.Now().Add(time.Second)

	jobUUID := uuid.New()
	testJob := Job{
		ID:          jobUUID,
		Title:       "Test Job",
		Description: "Test Job Description",
		Schedule:    "5s",
		Func:        nil,
		Status:      "",
		ScheduleId:  0,
	}

	r.register <- &testJob

	for time.Now().Before(deadline) {
		r.mutex.RLock()
		job, ok := r.jobs[jobUUID]
		r.mutex.RUnlock()

		if ok {
			if job.ScheduleId == 0 {
				t.Fatal("schedule id not assigned")
			}
			return
		}

		time.Sleep(time.Millisecond)
	}

	t.Fatal("job was never registered")
}

func TestUnregisterJob(t *testing.T) {
	r := New()

	testJob := Job{
		ID:       uuid.New(),
		Schedule: "5s",
		Func:     func() {},
	}

	// Register
	r.register <- &testJob

	// Wait until registered
	deadline := time.Now().Add(time.Second)

	for time.Now().Before(deadline) {
		r.mutex.RLock()
		_, ok := r.jobs[testJob.ID]
		r.mutex.RUnlock()

		if ok {
			break
		}

		time.Sleep(time.Millisecond)
	}

	// Unregister
	r.unregister <- testJob.ID

	// Wait until removed
	deadline = time.Now().Add(time.Second)

	for time.Now().Before(deadline) {
		r.mutex.RLock()
		_, ok := r.jobs[testJob.ID]
		r.mutex.RUnlock()

		if !ok {
			return
		}

		time.Sleep(time.Millisecond)
	}

	t.Fatal("job was not unregistered")
}
