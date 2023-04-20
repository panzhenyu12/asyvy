package worker

import (
	"asyvy/pkg/task"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type NodeConfig struct {
	Concurrency int
	RedisAddr   string
	RedisPwd    string
	CacheDir    string
}

type WorkNode struct {
	Conf        *NodeConfig
	asynqconfig *asynq.Config
}

func NewWorkNode(conf *NodeConfig) *WorkNode {
	return &WorkNode{
		Conf: conf,
	}
}

func (w *WorkNode) Run() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     w.Conf.RedisAddr,
			Password: w.Conf.RedisPwd,
		},
		asynq.Config{
			Concurrency: w.Conf.Concurrency,
			Logger:      logrus.StandardLogger(),
		},
	)
	mux := asynq.NewServeMux()
	handle := asynq.HandlerFunc(task.HandleImageScanTask)
	//mux.HandleFunc(task.TypeScanImage, task.HandleImageScanTask)
	mux.Handle(task.TypeScanImage, task.ScanImageMiddleware(handle))
	if err := srv.Run(mux); err != nil {
		logrus.Fatalf("could not run server: %v", err)
	}
}
