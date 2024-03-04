package executor

import "io"

// Executor interface
type Executor interface {
	RunCommand(image string, cmd []string) (io.ReadCloser, error)
}
