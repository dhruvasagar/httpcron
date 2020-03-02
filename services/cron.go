package services

import (
	"github.com/asdine/storm"
	"github.com/dhruvasagar/httpcron/db"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

type Cron struct {
	cron *cron.Cron
}

func (c *Cron) Get(entryID cron.EntryID) cron.Entry {
	return c.cron.Entry(entryID)
}

func (c *Cron) All() []cron.Entry {
	return c.cron.Entries()
}

func (c *Cron) Add(spec, command string) (cron.EntryID, error) {
	return c.cron.AddJob(spec, &Job{command: command})
}

func (c *Cron) Remove(entryID cron.EntryID) {
	c.cron.Remove(entryID)
}

func (c *Cron) Update(
	entryID cron.EntryID,
	spec, command string,
) (cron.EntryID, error) {
	c.cron.Remove(entryID)
	return c.Add(spec, command)
}

func loadJobsFromDB(cron *cron.Cron, sdb *storm.DB) error {
	var cronEntries []db.CronEntry
	err := sdb.All(&cronEntries)
	if err != nil {
		return err
	}

	for _, cronEntry := range cronEntries {
		entryID, err := cron.AddJob(cronEntry.Spec, &Job{
			command: cronEntry.Command,
		})
		if err != nil {
			log.Errorf("Failed to register db cronEntry: %v in cron, err: %s\n", cronEntry, err)
		} else {
			sdb.UpdateField(&cronEntry, "EntryID", entryID)
		}
	}

	return nil
}

func NewCron(sdb *storm.DB) (*Cron, error) {
	cron := cron.New()
	cron.Start()

	err := loadJobsFromDB(cron, sdb)
	if err != nil {
		return nil, err
	}

	return &Cron{
		cron: cron,
	}, nil
}
