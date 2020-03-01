package db

import "github.com/robfig/cron/v3"

type CronEntry struct {
	ID      int `storm:"id,increment"`
	Spec    string
	Command string
	EntryID cron.EntryID
}
