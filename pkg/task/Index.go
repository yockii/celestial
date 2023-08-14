package task

import (
	"context"
	"github.com/robfig/cron/v3"
	"time"
)

type taskPool struct {
	*cron.Cron
}
type schedule struct {
	duration time.Duration
}

func (s schedule) Next(t time.Time) time.Time {
	return t.Add(s.duration)
}

type job struct {
	f func()
}

func (j job) Run() {
	j.f()
}

var defaultTaskPool *taskPool

var parser = cron.NewParser(
	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)

func init() {
	defaultTaskPool = &taskPool{
		cron.New(cron.WithParser(parser)),
	}
}

func (t *taskPool) AddFunc(spec string, f func()) (int, error) {
	id, err := t.Cron.AddFunc(spec, f)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
func (t *taskPool) Remove(id int) {
	t.Cron.Remove(cron.EntryID(id))
}
func (t *taskPool) Schedule(duration time.Duration, f func()) int {
	id := t.Cron.Schedule(
		schedule{duration},
		job{f},
	)
	return int(id)
}
func (t *taskPool) Start() {
	t.Cron.Start()
}
func (t *taskPool) Stop() context.Context {
	return t.Cron.Stop()
}

func AddFunc(spec string, f func()) (int, error) {
	return defaultTaskPool.AddFunc(spec, f)
}
func Schedule(d time.Duration, f func()) int {
	return defaultTaskPool.Schedule(d, f)
}
func Remove(id int) {
	defaultTaskPool.Remove(id)
}
func Start() {
	defaultTaskPool.Start()
}
func Stop() context.Context {
	return defaultTaskPool.Stop()
}
