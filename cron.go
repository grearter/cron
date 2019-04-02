package cron

import (
	"errors"
	"log"
	"sort"
	"time"
)

type Entry struct {
	schedule *Schedule
	task     *Task
	next     time.Time
}

type Cron struct {
	entries []*Entry
	add     chan *Entry
	stop    chan struct{}
	running bool
	Logger  Logger
}

func NewCron() *Cron {
	return &Cron{
		entries: make([]*Entry, 0),
		add:     make(chan *Entry),
		stop:    make(chan struct{}),
		running: false,
	}
}

func (c *Cron) AddTask(name string, scheduleStr string, fn func(...interface{}) error, args ...interface{}) (err error) {

	if fn == nil {
		err = errors.New("fn can not be nil")
		return
	}

	schedule, err := Parse(scheduleStr)
	if err != nil {
		return
	}

	entry := &Entry{
		schedule: schedule,
		task: &Task{
			name: name,
			fn:   fn,
			args: args,
		},
	}

	if !c.running {
		c.entries = append(c.entries, entry)
	} else {
		c.add <- entry
	}

	return
}

func (c *Cron) IsRunning() bool {
	return c.running
}

func (c *Cron) Stop() {
	c.stop <- struct{}{}
}

func (c *Cron) Start() {
	if c.running {
		return
	}

	if c.Logger != nil {
		c.Logger.Info("%s cron start\n", time.Now().Format("2006-01-02 15:04:05"))
	} else {
		log.Printf("%s cron start\n", time.Now().Format("2006-01-02 15:04:05"))
	}
	c.running = true

	go func() {
		now := time.Now()

		// 计算任务的下一次执行时间
		for _, entry := range c.entries {
			entry.next = entry.schedule.Next(now)
		}

		for {
			// 按时间排序
			sort.Slice(c.entries, func(i, j int) bool {
				return c.entries[i].next.Before(c.entries[j].next)
			})

			// 设置timer
			var timer *time.Timer
			if len(c.entries) <= 0 || c.entries[0].next.IsZero() {
				timer = time.NewTimer(100000 * time.Hour)
			} else {
				timer = time.NewTimer(c.entries[0].next.Sub(now))
			}

			// select
			select {
			case now = <-timer.C:
				for _, entry := range c.entries {

					if entry.next.After(now) || entry.next.IsZero() {
						break
					}

					go c.runTask(entry.task)
					entry.next = entry.schedule.Next(now)
				}

			case newEntry := <-c.add:
				timer.Stop()
				now = time.Now()
				newEntry.next = newEntry.schedule.Next(now)
				c.entries = append(c.entries, newEntry)

			case <-c.stop:
				timer.Stop()
				return
			}
		}
	}()

	return
}

func (c *Cron) runTask(task *Task) {

	if c.Logger != nil {
		c.Logger.Info("%s task:%s begin\n", task.name, time.Now().Format("2006-01-02 15:04:05"))
	} else {
		log.Printf("%s task:%s begin\n", task.name, time.Now().Format("2006-01-02 15:04:05"))
	}

	err := task.fn(task.args...)

	if c.Logger != nil {
		if err != nil {
			c.Logger.Error("%s task:%s end, err: %v\n", task.name, time.Now().Format("2006-01-02 15:04:05"), err)
		} else {
			c.Logger.Info("%s task:%s end\n", task.name, time.Now().Format("2006-01-02 15:04:05"))
		}
	} else {
		if err != nil {
			log.Printf("%s task:%s end, err: %v\n", task.name, time.Now().Format("2006-01-02 15:04:05"), err)
		} else {
			log.Printf("%s task:%s end\n", task.name, time.Now().Format("2006-01-02 15:04:05"))
		}
	}

	return
}
