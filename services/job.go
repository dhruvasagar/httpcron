package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func getJobTimeout() time.Duration {
	jobTimeoutStr := os.Getenv("JOB_TIMEOUT")
	jobTimeoutInt, _ := strconv.Atoi(jobTimeoutStr)
	jobTimeout := time.Duration(jobTimeoutInt)
	if jobTimeout == 0 {
		jobTimeout = 100 * time.Millisecond
	}
	return jobTimeout
}

type Job struct {
	command string
}

// Implement cron.Job interface
func (j *Job) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), getJobTimeout())
	defer cancel()

	cmd := exec.CommandContext(ctx, j.command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Error while executing command", err)
		return
	}
	fmt.Println(output)
}
