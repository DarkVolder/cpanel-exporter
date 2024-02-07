package scheduler

import (
	"cpanel_exporter/metrics"
	"github.com/go-co-op/gocron"
	"time"
)

type scheduler struct {
	interval      int
	intervalHeavy int
	sched         *gocron.Scheduler
	metrics       metrics.Metrics
}

func New(interval int, intervalHeavy int, metrics metrics.Metrics) *scheduler {
	s := &scheduler{
		interval:      interval,
		intervalHeavy: intervalHeavy,
		metrics:       metrics,
		sched:         gocron.NewScheduler(time.UTC),
	}
	s.SetTasks()
	return s
}

func (s *scheduler) Run() {
	s.sched.StartBlocking()
}

func (s *scheduler) SetTasks() {
	s.sched.Every(s.interval).Seconds().Do(
		func() {
			s.metrics.FetchMetrics()
		},
	)
	s.sched.Every(s.intervalHeavy).Seconds().Do(
		func() {
			s.metrics.FetchUpiMetrics()
		},
	)
}
