package executor

import (
	"context"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// DockerExecutor implements the executor interface for Docker
type DockerExecutor struct {
	cli *client.Client
}

// NewDockerExecutor instance creator
func NewDockerExecutor() (*DockerExecutor, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerExecutor{cli: cli}, nil
}

// Run command inside a container and return  the logs
func (de *DockerExecutor) RunCommand(image string, cmd []string) error {
	ctx := context.Background()
	resp, err := de.cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   cmd,
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return err
	}

	if err := de.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := de.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := de.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	err = de.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})
	if err != nil {
		panic(err)
	}
}
