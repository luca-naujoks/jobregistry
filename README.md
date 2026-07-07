# JobRegistry

A Go module for registering and managing jobs.

`JobRegistry` provides a structured way to register, manage, and execute jobs within your Go applications. The Registry keeps track and manages your jobs while the [robfig/cron/v3](https://github.com/robfig/cron) package is handling the Execution on Time.

## Getting Started

To use this module in your Go project:

```bash
go get github.com/luca-naujoks/jobregistry
```

## Usage

Import the package into your Go code:

 ### 1. Create a registry

```go
r := registry.New()
```

### 2. Define a job

```go
jobUUID := uuid.New()
testJob := registry.Job{
    ID:          jobUUID,
    Title:       "Test Job",
    Description: "Test Job Description",
    Schedule:    "5s",
    Func:        nil,
    Status:      "",
    ScheduleId:  0,
}
```

### 3. Register the job

```go
err := r.Register(testJob)
if err != nil {
    return
}
```

### 4. Read the last run time

```go
jobNextRun, err := r.LastRun(jobUUID)
if err != nil {
    return
}

fmt.Println(jobNextRun.Format("2006-01-02 15:04"))
```
