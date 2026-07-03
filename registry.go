package registry

import (
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

func New() *Registry {
	r := &Registry{
		jobs:      make(map[uuid.UUID]*Job),
		scheduler: cron.New(),
	}
	r.scheduler.Start()

	return r
}
