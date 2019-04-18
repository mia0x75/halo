package crons

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
)

// Task holds information about task
type Task struct {
	Schedule
	Name      string
	UUID      string
	Cmd       string
	Status    string
	IsRunning bool
	Params    []string
	sync.Mutex
}

// NewTask returns an instance of task
func NewTask(name, cmd string, params []string) *Task {
	return &Task{
		Name:      name,
		Cmd:       cmd,
		Params:    params,
		IsRunning: false,
		Status:    "P",
	}
}

// NewTaskWithSchedule creates an instance of task with the provided schedule information
func NewTaskWithSchedule(name, cmd string, params []string, schedule Schedule) *Task {
	return &Task{
		Name:      name,
		Schedule:  schedule,
		Cmd:       cmd,
		Params:    params,
		IsRunning: false,
		Status:    "P",
	}
}

// IsDue returns a boolean indicating whether the task should execute or not
func (t *Task) IsDue() bool {
	timeNow := time.Now()
	t.Lock()
	defer t.Unlock()
	return timeNow == t.NextRun || timeNow.After(t.NextRun)
}

// Run will execute the task and schedule it's next run.
func (t *Task) Run() {
	t.Lock()
	defer t.Unlock()
	if t.IsRunning {
		return
	}
	t.IsRunning = true
	defer func() {
		t.IsRunning = false
	}()

	if _, err := tools.TimeoutedExec(math.MaxInt32*time.Second, t.Cmd, t.Params...); err != nil {
		fmt.Println(err)
		t.Status = "F"
	} else {
		t.Status = "S"
	}

	t.LastRun = time.Now()

	if !t.IsRecurring {
		t.NextRun = time.Time{}
		return
	}
	t.NextRun = t.NextRun.Add(t.Interval)
}

// ToCron 将task结构转成models.Cron
func (t *Task) ToCron() *models.Cron {
	t.Lock()
	defer t.Unlock()
	cron := &models.Cron{
		UUID:      t.UUID,
		Status:    t.Status,
		Name:      t.Name,
		Cmd:       t.Cmd,
		NextRun:   t.NextRun.Format(time.RFC3339),
		LastRun:   t.LastRun.Format(time.RFC3339),
		Recurrent: 0,
		Interval:  t.Interval.String(),
	}
	if t.IsRecurring {
		cron.Recurrent = 1
	}
	var data []byte
	data, _ = json.Marshal(t.Params)
	cron.Params = string(data)
	return cron
}
