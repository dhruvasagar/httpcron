package services

import "github.com/robfig/cron/v3"

type Cron struct {
	cron *cron.Cron
}

func (c *Cron) Add(spec, command string) (cron.EntryID, error) {
	return c.cron.AddJob(spec, &Job{command: command})
}

func (c *Cron) Remove(entryID cron.EntryID) {
	c.cron.Remove(entryID)
}

func (c *Cron) Update(entryID cron.EntryID, spec, command string) (cron.EntryID, error) {
	c.cron.Remove(entryID)
	return c.Add(spec, command)
}

func NewCron() *Cron {
	cron := cron.New()
	cron.Start()

	return &Cron{
		cron: cron,
	}
}
