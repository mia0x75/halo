package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/service"
)

// 1、任务完成或者失败，持久化到数据库

var (
	cmd       = "halocli"
	scheduler *Scheduler
	once      sync.Once
)

// reloadInterval is how often the scheduler re-reads the schedulable crons.
const reloadInterval = 60 * time.Second

// repos is injected once at startup via Init so the scheduler reloads and
// persists tasks through the repository layer.

var repos service.Repos

// Init injects the repositories the scheduler reads and writes crons with. It
// must run before NewScheduler.
func Init(r service.Repos) {
	repos = r
}

// Schedule 保存任务执行时间相关的信息
type Schedule struct {
	IsRecurring bool
	LastRun     time.Time
	NextRun     time.Time
	Interval    time.Duration
}

// Scheduler 用于调度任务，它保存了待执行任务的相关信息
type Scheduler struct {
	sync.Mutex
	stopChan chan bool
	tasks    map[string]*Task
}

// NewScheduler 返回一个Scheduler结构体的实例
func NewScheduler() *Scheduler {
	once.Do(func() {
		scheduler = &Scheduler{
			stopChan: make(chan bool),
			tasks:    make(map[string]*Task),
		}
		scheduler.Start()
	})
	return scheduler
}

// RunAt 在一个给定时间执行一个任务
func (s *Scheduler) RunAt(when time.Time, name string, params ...string) (string, error) {
	task := NewTask(name, cmd, params)
	task.NextRun = when
	task.LastRun, _ = time.Parse(time.RFC3339, "2006-01-02 15:04:05")
	s.register(task)
	return task.UUID, nil
}

// RunAfter 等待一个指定的时间后执行一个任务
func (s *Scheduler) RunAfter(duration time.Duration, name string, params ...string) (string, error) {
	return s.RunAt(time.Now().Add(duration), name, params...)
}

// RunEvery 给定周期，循环执行某一个任务
func (s *Scheduler) RunEvery(interval time.Duration, name string, params ...string) (string, error) {
	t, exists := s.taskExists(cmd)

	if exists && t.Interval == interval {
		return t.UUID, nil
	}

	task := NewTask(name, cmd, params)
	task.IsRecurring = true
	task.Interval = interval
	task.LastRun, _ = time.Parse(time.RFC3339, "2006-01-02 15:04:05")
	if exists {
		task.NextRun = t.LastRun.Add(interval)
	} else {
		task.NextRun = time.Now().Add(interval)
	}

	s.register(task)
	return task.UUID, nil
}

// Start 启动任务调度程序
func (s *Scheduler) Start() error {
	if err := s.load(); err != nil {
		return err
	}
	s.runPending()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		reload := time.NewTicker(reloadInterval)
		for {
			select {
			case <-ticker.C:
				s.runPending()
			case <-reload.C:
				if err := s.reload(); err != nil {
					log.Errorf("[E] Failed to reload scheduled tasks, err: %s", err.Error())
				}
			case <-s.stopChan:
				close(s.stopChan)
			}
		}
	}()

	return nil
}

// Stop will put the scheduler to halt
func (s *Scheduler) Stop() {
	s.stopChan <- true
}

// Wait is a convenience function for blocking until the scheduler is stopped.
func (s *Scheduler) Wait() {
	<-s.stopChan
}

// Cancel is used to cancel the planned execution of a specific task using it's ID.
// The ID is returned when the task was scheduled using RunAt, RunAfter or RunEvery
func (s *Scheduler) Cancel(uuid string) error {
	s.Lock()
	task, found := s.tasks[uuid]
	s.Unlock()
	if !found {
		return fmt.Errorf("Task not found")
	}

	// The stored row has to leave the schedulable set, otherwise the periodic
	// reload picks the cancelled schedule up again after a restart.
	task.Lock()
	task.Status = statusCancelled
	task.Unlock()

	s.persist(task)
	s.Lock()
	defer s.Unlock()
	delete(s.tasks, uuid)
	return nil
}

// Clear 取消所有的调度任务
func (s *Scheduler) Clear() {
	s.Lock()
	defer s.Unlock()
	for uuid, task := range s.tasks {
		s.persist(task)
		delete(s.tasks, uuid)
	}
}

/*
P - Pending 等待执行
C - Cancelled - 任务取消
S - Succeed - 执行成功
F - Failure - 执行失败
R - Runing - 正在执行
H - Halt - 任务暂停
*/
func (s *Scheduler) load() error {
	tasks, err := s.schedulableTasks()
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()
	s.tasks = tasks
	return nil
}

// schedulableTasks builds the task set from the schedulable crons.
func (s *Scheduler) schedulableTasks() (map[string]*Task, error) {
	domainCrons, err := repos.Cron.ListSchedulable(context.Background())
	if err != nil {
		return nil, err
	}
	tasks := map[string]*Task{}
	for _, c := range domainCrons {
		cron := bridge.CronFromDomain(c)
		task, err := taskFromCron(cron)
		if err != nil {
			return nil, err
		}
		tasks[cron.UUID] = task
	}
	return tasks, nil
}

// taskFromCron rebuilds an in-memory task from its stored schedule.
func taskFromCron(cron *models.Cron) (*Task, error) {
	lastRun, err := time.Parse(time.RFC3339, cron.LastRun)
	if err != nil {
		return nil, fmt.Errorf("cron %s last_run: %w", cron.UUID, err)
	}

	nextRun, err := time.Parse(time.RFC3339, cron.NextRun)
	if err != nil {
		return nil, fmt.Errorf("cron %s next_run: %w", cron.UUID, err)
	}

	interval, err := time.ParseDuration(cron.Interval)
	if err != nil {
		return nil, fmt.Errorf("cron %s interval: %w", cron.UUID, err)
	}

	params := []string{}
	if err := json.Unmarshal([]byte(cron.Params), &params); err != nil {
		return nil, fmt.Errorf("cron %s params: %w", cron.UUID, err)
	}

	task := NewTaskWithSchedule(cron.Name, cron.Cmd, params, Schedule{
		IsRecurring: cron.Recurrent == 1,
		Interval:    interval,
		LastRun:     lastRun,
		NextRun:     nextRun,
	})
	task.UUID = cron.UUID
	return task, nil
}

// reload refreshes the task set from the database: crons created after startup
// are picked up, and crons deleted in the meantime are dropped. Tasks that are
// running keep their in-memory state.
func (s *Scheduler) reload() error {
	tasks, err := s.schedulableTasks()
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()
	for uuid, task := range tasks {
		if _, ok := s.tasks[uuid]; !ok {
			s.tasks[uuid] = task
		}
	}
	for uuid, task := range s.tasks {
		if _, ok := tasks[uuid]; !ok && !task.IsRunning {
			delete(s.tasks, uuid)
		}
	}
	return nil
}

// Reload refreshes the running scheduler from the database. It is called after
// the cron rows change so a new schedule does not wait for a restart.
func Reload() {
	if scheduler == nil {
		return
	}
	if err := scheduler.reload(); err != nil {
		log.Errorf("[E] Failed to reload scheduled tasks, err: %s", err.Error())
	}
}

func (s *Scheduler) runPending() {
	s.Lock()
	defer s.Unlock()
	wg := &sync.WaitGroup{}
	for _, task := range s.tasks {
		if task.IsDue() {
			wg.Add(1)
			go s.runTask(task, wg)
		}
	}
	wg.Wait()
	for _, task := range s.tasks {
		if !task.IsRecurring && task.NextRun.Before(time.Now()) {
			s.removeTask(task)
		}
	}
}

func (s *Scheduler) runTask(task *Task, wg *sync.WaitGroup) {
	defer wg.Done()
	task.Run()
	task.Lock()
	defer task.Unlock()
	s.persist(task)
}

func (s *Scheduler) removeTask(task *Task) {
	delete(s.tasks, task.UUID)
}

func (s *Scheduler) register(task *Task) {
	err := s.persist(task)
	if err != nil {
		log.Errorf("[E] Failed to persist task, err: %s", err.Error())
	}
	s.Lock()
	s.tasks[task.UUID] = task
	s.Unlock()
}

func (s *Scheduler) taskExists(name string) (*Task, bool) {
	for _, t := range s.tasks {
		if t.Cmd == name {
			return t, true
		}
	}
	return nil, false
}

func (s *Scheduler) persist(task *Task) (err error) {
	ctx := context.Background()
	cron := &models.Cron{}
	if task.UUID != "" {
		dc, getErr := repos.Cron.GetByUUID(ctx, task.UUID)
		if getErr != nil {
			log.Errorf("[E] An unexpected error ocurred, err: %s", getErr.Error())
			return getErr
		}
		cron = bridge.CronFromDomain(dc)
		cron.Status = task.Status
		cron.LastRun = task.LastRun.Format(time.RFC3339)
		cron.NextRun = task.NextRun.Format(time.RFC3339)
		cron.BeforeUpdate()
		if err = repos.Cron.Update(ctx, bridge.CronToDomain(cron)); err != nil {
			log.Errorf("[E] An unexpected error ocurred, err: %s", err.Error())
			return err
		}
	} else {
		cron.Name = task.Name
		cron.Cmd = task.Cmd
		cron.NextRun = task.NextRun.Format(time.RFC3339)
		cron.LastRun = task.LastRun.Format(time.RFC3339)
		cron.Recurrent = 0
		cron.Status = task.Status
		if task.IsRecurring {
			cron.Recurrent = 1
		}
		cron.Interval = task.Interval.String()

		var data []byte
		data, err = json.Marshal(task.Params)
		cron.Params = string(data)

		cron.BeforeInsert()
		newID, createErr := repos.Cron.Create(ctx, bridge.CronToDomain(cron))
		if createErr != nil {
			return createErr
		}
		cron.CronID = uint(newID)
		task.UUID = cron.UUID
	}

	// 任务执行记录：调度即建行，执行后回写状态与失败原因。
	if err = s.syncTask(task, cron.CronID); err != nil {
		log.Errorf("[E] Failed to record task execution, err: %s", err.Error())
	}
	return
}

// JobCrons 返回队列中的系统作业列表（不含工单预约），并转化成models.Cron
func (s *Scheduler) JobCrons() (L []*models.Cron, err error) {
	s.Lock()
	defer s.Unlock()
	for _, task := range s.tasks {
		if task.IsTicketBooking() {
			continue
		}
		L = append(L, task.ToCron())
	}
	return
}
