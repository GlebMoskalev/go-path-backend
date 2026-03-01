package service

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/GlebMoskalev/go-path-backend/internal/model"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"go.uber.org/zap"
)

type SandboxService struct {
	log     *zap.Logger
	docker  *client.Client
	image   string
	timeOut time.Duration
	memory  int64
}

func NewSandboxService(log *zap.Logger, imageName string, timeOut time.Duration, memory int64) (*SandboxService, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("docker client: %v", err)
	}

	_, err = docker.ImageInspect(context.Background(), imageName)
	if err != nil {
		return nil, fmt.Errorf("sandbox image %s not found locally: %v", imageName, err)
	}

	return &SandboxService{
		log:     log,
		docker:  docker,
		image:   imageName,
		timeOut: timeOut,
		memory:  memory,
	}, nil
}

func (s *SandboxService) Run(ctx context.Context, userCode, testFile string) model.SubmitResult {
	files := map[string]string{
		"go.mod":           "module solution\n\ngo 1.25\n",
		"solution.go":      userCode,
		"solution_test.go": testFile,
	}

	tarBuf, err := s.createTarArchive(files)
	if err != nil {
		s.log.Error("failed to create tar archive", zap.Error(err))
		return model.SubmitResult{Error: "internal error"}
	}

	ctx, cancel := context.WithTimeout(ctx, s.timeOut)
	defer cancel()

	resp, err := s.docker.ContainerCreate(ctx, &container.Config{
		Image:      s.image,
		Cmd:        []string{"go", "test", "-v", "-json", "-count=1", "./..."},
		WorkingDir: "/sandbox",
	}, &container.HostConfig{
		NetworkMode: "none",
		Resources: container.Resources{
			Memory:   s.memory,
			NanoCPUs: 500000000,
		},
	}, nil, nil, "")
	if err != nil {
		s.log.Error("container create failed", zap.Error(err))
		return model.SubmitResult{Error: "internal error"}
	}

	defer s.docker.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})

	if err := s.docker.CopyToContainer(ctx, resp.ID, "/sandbox", tarBuf, container.CopyToContainerOptions{}); err != nil {
		s.log.Error("copy to container failed", zap.Error(err))
		return model.SubmitResult{Error: "internal error"}
	}

	if err := s.docker.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		s.log.Error("container start failed", zap.Error(err))
		return model.SubmitResult{Error: "internal error"}
	}

	statusCh, errCh := s.docker.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			s.log.Error("container wait failed", zap.Error(err))
			return model.SubmitResult{Error: "execution timeout"}
		}
	case <-statusCh:
	}

	out, err := s.docker.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		s.log.Error("container logs failed", zap.Error(err))
		return model.SubmitResult{Error: "internal error"}
	}
	defer out.Close()

	var stdout, stderr bytes.Buffer
	stdcopy.StdCopy(&stdout, &stderr, out)

	return s.parseTestOutput(stdout.String(), stderr.String())
}

type goTestEvent struct {
	Action  string  `json:"Action"`
	Test    string  `json:"Test"`
	Output  string  `json:"Output"`
	Elapsed float64 `json:"Elapsed"`
}

func (s *SandboxService) parseTestOutput(stdout, stderr string) model.SubmitResult {
	if stderr != "" && !strings.Contains(stdout, `"Action"`) {
		return model.SubmitResult{
			Passed: false,
			Error:  stderr,
		}
	}

	lines := strings.Split(strings.TrimSpace(stdout), "\n")

	testOutputs := make(map[string]string)
	testResults := make(map[string]bool)

	for _, line := range lines {
		var event goTestEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}

		if event.Test == "" {
			continue
		}

		switch event.Action {
		case "output":
			testOutputs[event.Test] += event.Output
		case "pass":
			testResults[event.Test] = true
		case "fail":
			testResults[event.Test] = false
		}
	}

	if len(testResults) == 0 {
		return model.SubmitResult{
			Passed: false,
			Error:  stdout + stderr,
		}
	}

	allPassed := true
	var tests []model.TestResult

	for name, passed := range testResults {
		if !passed {
			allPassed = false
		}
		tests = append(tests, model.TestResult{
			Name:   name,
			Passed: passed,
			Output: testOutputs[name],
		})
	}

	return model.SubmitResult{
		Passed: allPassed,
		Tests:  tests,
	}
}

func (s *SandboxService) createTarArchive(files map[string]string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for name, content := range files {
		hdr := &tar.Header{
			Name: name,
			Mode: 0644,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, fmt.Errorf("tar write header %s: %w", name, err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			return nil, fmt.Errorf("tar write %s: %w", name, err)
		}
	}
	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("tar close: %w", err)
	}
	return &buf, nil
}
