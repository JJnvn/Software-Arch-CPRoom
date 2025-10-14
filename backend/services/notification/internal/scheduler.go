package internal

import (
	"context"
	"log"
	"time"
)

type Scheduler struct {
	service  *NotificationService
	interval time.Duration
}

func NewScheduler(service *NotificationService, interval time.Duration) *Scheduler {
	if interval <= 0 {
		interval = time.Minute
	}
	return &Scheduler{
		service:  service,
		interval: interval,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.service.ProcessDueNotifications(ctx); err != nil {
				log.Printf("scheduler error: %v", err)
			}
		}
	}
}
