package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// Add bash scripts pull golang image before executing executor
func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "golang",
		Cmd:   []string{"go", "adada"},
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		log.Fatal(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Fatal(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatal(err)

		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		log.Fatal(err)

	}

	fmt.Println("########################")
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	err = cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})
	if err != nil {
		log.Fatal(err)

	}
}
