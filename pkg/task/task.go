package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

const (
	TypeScanImage = "scan:docker"
)

type ImageScanTaskPayload struct {
	// Image is the image to scan
	Image string `json:"image"`
	// RegistryUsername is the username to use when scanning the registry
	RegistryUsername string `json:"registry_username"`
	// RegistryPassword is the password to use when scanning the registry
	RegistryPassword string `json:"registry_password"`
}

func ImageScanTask(image, username, password string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageScanTaskPayload{
		Image:            image,
		RegistryUsername: username,
		RegistryPassword: password,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeScanImage, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute), asynq.Retention(2*24*time.Hour)), nil
}

func HandleImageScanTask(ctx context.Context, t *asynq.Task) error {
	var payload ImageScanTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return ScanImage(payload.Image, payload.RegistryUsername, payload.RegistryPassword)
}

func ScanImage(image, username, password string) error {
	// ...
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("trivy", "image", "--image-config-scanners", "config", "--format", "json", "--no-progress", "--cache-dir",
		"/Users/panzhenyu/tryvy-db", "--skip-db-update", "--skip-java-db-update", image)
	if len(username) > 0 && len(password) > 0 {
		cmd.Env = append(cmd.Env, "TRIVY_USERNAME="+username)
		cmd.Env = append(cmd.Env, "TRIVY_PASSWORD="+password)
	}
	//cmd.Env = append(cmd.Env, "TRIVY_USERNAME="+username)
	//cmd.Env = append(cmd.Env, "TRIVY_PASSWORD="+password)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	logrus.Infof("Scan result: %s", out.String())
	return nil
}

func ScanImageMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		start := time.Now()
		logrus.Printf("Start processing %q", t.Type())
		err := h.ProcessTask(ctx, t)
		if err != nil {
			return err
		}
		logrus.Printf("Finished processing %q: Elapsed Time = %v", t.Type(), time.Since(start))
		return nil
	})
}
