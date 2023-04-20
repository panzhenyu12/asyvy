package controller

import "github.com/hibiken/asynq"

type Controller struct {
	TaskClient *asynq.Client
}
