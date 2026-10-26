package worker

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/tools"
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
	// Err keeps the failure reason of the last run so it can be persisted
	// with the task execution.
	Err error
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

	if _, err := tools.TimeoutedExec(math.MaxInt32*time.Second, resolveCommand(t.Cmd), t.Params...); err != nil {
		fmt.Println(err)
		t.Status = "F"
		t.Err = err
	} else {
		t.Status = "S"
		t.Err = nil
	}

	t.LastRun = time.Now()

	if !t.IsRecurring {
		t.NextRun = time.Time{}
		return
	}
	t.NextRun = t.NextRun.Add(t.Interval)
}

// IsTicketBooking reports whether the task is a one-off ticket execution
// booked via ScheduleTicket; those belong to the execution plan page.
func (t *Task) IsTicketBooking() bool {
	t.Lock()
	defer t.Unlock()
	return len(t.Params) > 0 && t.Params[0] == domain.TicketExecuteCmd
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

// resolveCommand returns the command to run. The bare CLI name is looked up
// next to the running binary first: a deployment keeps the server and the CLI
// together, and relying on PATH alone would break the scheduled jobs.
func resolveCommand(command string) string {
	if command != cmd || strings.ContainsRune(command, os.PathSeparator) {
		return command
	}
	executable, err := os.Executable()
	if err != nil {
		return command
	}
	candidate := filepath.Join(filepath.Dir(executable), cmd)
	if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
		return candidate
	}
	return command
}
