package executor

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"
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

// RunCommand executes a command inside a Docker container and returns the logs
func (de *DockerExecutor) RunCommand(ctx context.Context image string, cmd []string) error {

	// Create container
	resp, err := de.cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   cmd,
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	// Start container
	if err := de.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	// Periodically check container status
	go de.monitorContainerStatus(ctx, resp.ID)

	// Wait for container to finish
	statusCh, errCh := de.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error while waiting for container: %v", err)
		}
	case <-statusCh:
	}

	// Retrieve container logs
	out, err := de.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return fmt.Errorf("failed to retrieve container logs: %v", err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	// Remove container
	err = de.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove container: %v", err)
	}

	return nil
}

// monitorContainerStatus periodically checks the status of the container
func (de *DockerExecutor) monitorContainerStatus(ctx context.Context, containerID string) {
	ticker := time.NewTicker(10 * time.Second) // Adjust the interval as needed

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			containerInfo, err := de.cli.ContainerInspect(ctx, containerID)
			if err != nil {
				fmt.Printf("Error checking container status: %v\n", err)
				return
			}

			if containerInfo.State.Status == "running" {
				fmt.Printf("Container %s finished with status: %s\n", containerID, containerInfo.State.Status)
				continue
			}

			fmt.Printf("Container %s finished with status: %s\n", containerID, containerInfo.State.Status)

			if containerInfo.State.Dead{
				continue
			}
			return
		}
	}
}
