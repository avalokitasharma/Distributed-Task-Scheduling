package service

import (
	"errors"
	"time"

	"github.com/avalokitasharma/job-scheduler/scheduler-service/repository"
	"github.com/robfig/cron/v3"
)

func computeNextRun(job *repository.Job) (time.Time, error) {
	now := time.Now()

	switch job.ScheduleType {
	case "ONCE":
		return job.NextRunAt, nil

	case "CRON":
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, err := parser.Parse(job.CronExpression)
		if err != nil {
			return time.Time{}, err
		}
		return schedule.Next(now), nil

	case "INTERVAL":
		return now.Add(time.Duration(job.IntervalSeconds) * time.Second), nil

	default:
		return time.Time{}, errors.New("invalid schedule type")
	}
}
