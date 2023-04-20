package main

import (
	"asyvy/pkg/task"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "", Password: ""})
	defer client.Close()
	task, err := task.ImageScanTask("ubuntu", "", "")
	if err != nil {
		logrus.Error(err)
	}
	taskinfo, err := client.Enqueue(task)
	logrus.Println(taskinfo.ID)
}
