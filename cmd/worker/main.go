package main

import (
	"asyvy/config"
	"asyvy/pkg/worker"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	var confpath string
	app := &cli.App{
		Name:  "asyvy-worker",
		Usage: "asyvy-worker",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"f"},
				Value:       "",
				Usage:       "config file path",
				Destination: &confpath,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "trivy-worker",
				Aliases: []string{"a"},
				Usage:   "开启trivy扫描worker",
				Action: func(cCtx *cli.Context) error {
					fmt.Println(confpath)
					conf, err := config.GetConfig(confpath)
					if err != nil {
						return err
					}
					logrus.Info("init config success", conf)
					// 初始化log组件
					logrus.SetReportCaller(true)
					logrus.SetFormatter(&MyJSONFormatter{})
					logrus.Info("init log success")
					//logrus.SetFormatter(&MyJSONFormatter{})
					wrk := worker.NewWorkNode(&worker.NodeConfig{
						Concurrency: conf.Worker.Concurrency,
						CacheDir:    conf.Worker.CacheDir,
						RedisAddr:   conf.Redis.Addr,
						RedisPwd:    conf.Redis.Password,
					})
					wrk.Run()
					// 等到停止信号
					//AwaitSignal()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal(err.Error())
	}
}
func AwaitSignal() os.Signal {
	return WaitForSignal(syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
}
func WaitForSignal(sources ...os.Signal) os.Signal {
	var s = make(chan os.Signal, 1)
	defer signal.Stop(s) //the second Ctrl+C will force shutdown

	signal.Notify(s, sources...)
	return <-s //blocked
}

type MyJSONFormatter struct {
	logrus.JSONFormatter
}

func (f *MyJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Note this doesn't include Time, Level and Message which are available on
	// the Entry. Consult `godoc` on information about those fields or read the
	// source of the official loggers.
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			entry.Data[k] = fmt.Sprintf("%+v", v)
		}
	}
	return f.JSONFormatter.Format(entry)
}
