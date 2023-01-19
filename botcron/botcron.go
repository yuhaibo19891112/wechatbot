package botcron

import (
	"github.com/robfig/cron/v3"
)

func NewWeChatBotCron(cronStr string, cmd func()) *cron.Cron {
	crontab := cron.New(cron.WithSeconds())

	crontab.AddFunc(cronStr, cmd)

	crontab.Start()

	return crontab
}

func NewWeChatBotCronJob(cronStr string, job cron.Job) *cron.Cron {
	crontab := cron.New(cron.WithSeconds())

	crontab.AddJob(cronStr, job)

	crontab.Start()

	return crontab
}
