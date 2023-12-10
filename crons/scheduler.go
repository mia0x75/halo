package crons

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"xorm.io/builder"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// 1、任务完成或者失败，持久化到数据库

var (
	cmd       = "halocli"
	scheduler *Scheduler
	once      sync.Once
)

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
		for {
			select {
			case <-ticker.C:
				s.runPending()
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
	crons := []*models.Cron{}
	where := builder.Or(builder.Eq{"recurrent": 1}, builder.And(builder.Neq{"recurrent": 1}, builder.In("status", "P", "H")))
	if err := g.Engine.Where(where).Find(&crons); err != nil {
		return err
	}
	tasks := map[string]*Task{}
	s.Lock()
	defer s.Unlock()
	for _, cron := range crons {
		lastRun, err := time.Parse(time.RFC3339, cron.LastRun)
		if err != nil {
			log.Errorf("[E] %s", err.Error())
			return err
		}

		nextRun, err := time.Parse(time.RFC3339, cron.NextRun)
		if err != nil {
			log.Errorf("[E] %s", err.Error())
			return err
		}

		interval, err := time.ParseDuration(cron.Interval)
		if err != nil {
			log.Errorf("[E] %s", err.Error())
			return err
		}

		isRecurring := cron.Recurrent == 1

		params := []string{}
		err = json.Unmarshal([]byte(cron.Params), &params)
		if err != nil {
			log.Errorf("[E] %s", err.Error())
			return err
		}

		task := NewTaskWithSchedule(cron.Name, cron.Cmd, params, Schedule{
			IsRecurring: isRecurring,
			Interval:    time.Duration(interval),
			LastRun:     lastRun,
			NextRun:     nextRun,
		})
		task.UUID = cron.UUID
		tasks[cron.UUID] = task
	}

	s.tasks = tasks
	return nil
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
	cron := &models.Cron{}
	if task.UUID != "" {
		ok := false
		if ok, err = g.Engine.Where("`uuid` = ?", task.UUID).Get(cron); err != nil {
			log.Errorf("[E] An unexpected error ocurred, err: %s", err.Error())
			return err
		} else if !ok {
			return fmt.Errorf("[E] Cannot find cron (uuid=%s)", task.UUID)
		}
		cron.Status = task.Status
		cron.LastRun = task.LastRun.Format(time.RFC3339)
		cron.NextRun = task.NextRun.Format(time.RFC3339)
		if _, err = g.Engine.ID(cron.CronID).Update(cron); err != nil {
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

		if _, err = g.Engine.Insert(cron); err != nil {
			return err
		}

		task.UUID = cron.UUID
	}
	return
}

// Crons 返回队列中的任务列表，并转化成models.Cron
func (s *Scheduler) Crons() (L []*models.Cron, err error) {
	s.Lock()
	defer s.Unlock()
	for _, task := range s.tasks {
		L = append(L, task.ToCron())
	}
	return
}
