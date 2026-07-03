package registry

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJobReturnsExistingJob(t *testing.T) {
	r := New()

	job := Job{
		ID:          uuid.New(),
		Title:       "test",
		Description: "",
		Schedule:    "1s",
		Func:        func() {},
		Status:      "",
		ScheduleId:  0,
	}

	if err := r.Register(job); err != nil {
		t.Fatal(err)
	}

	time.Sleep(20 * time.Millisecond)

	got, err := r.Job(job.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ID != job.ID {
		t.Errorf("expected %s, got %s", job.ID, got.ID)
	}
}

func TestJobReturnsErrorWhenMissing(t *testing.T) {
	r := New()

	_, err := r.Job(uuid.New())

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestJobsReturnsAllJobs(t *testing.T) {
	r := New()

	j1 := Job{
		ID:       uuid.New(),
		Title:    "job1",
		Schedule: "1s",
		Func:     func() {},
	}

	j2 := Job{
		ID:       uuid.New(),
		Title:    "job2",
		Schedule: "1s",
		Func:     func() {},
	}

	r.Register(j1)
	r.Register(j2)

	time.Sleep(20 * time.Millisecond)

	jobs := r.Jobs()

	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
	}
}

func TestRemoveDeletesJob(t *testing.T) {
	r := New()

	job := Job{
		ID:       uuid.New(),
		Title:    "test",
		Schedule: "1s",
		Func:     func() {},
	}

	r.Register(job)

	time.Sleep(20 * time.Millisecond)

	r.Remove(job.ID)

	time.Sleep(20 * time.Millisecond)

	_, err := r.Job(job.ID)
	if err == nil {
		t.Fatal("expected job to be removed")
	}
}

func TestRunExecutesJob(t *testing.T) {
	r := New()

	executed := make(chan struct{})

	job := Job{
		ID:       uuid.New(),
		Title:    "run",
		Schedule: "1h",
		Func: func() {
			close(executed)
		},
	}

	r.Register(job)

	time.Sleep(20 * time.Millisecond)

	if err := r.Run(job.ID); err != nil {
		t.Fatal(err)
	}

	select {
	case <-executed:
	case <-time.After(time.Second):
		t.Fatal("job was not executed")
	}
}

func TestRunMissingJob(t *testing.T) {
	r := New()

	err := r.Run(uuid.New())
	if err == nil {
		t.Fatal("expected error")
	}
}
