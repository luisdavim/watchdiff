//go:build !windows
// +build !windows

package watcher

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"syscall"
)

const DefaultShell = "sh"

func execute(ctx context.Context, cmd string, args []string, captureStderr bool) ([]byte, int) {
	var stdoutBuf bytes.Buffer
	c := exec.CommandContext(ctx, cmd, args...)
	c.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		// Pdeathsig: syscall.SIGKILL,
	}

	c.Stdout = &stdoutBuf
	if captureStderr {
		c.Stderr = &stdoutBuf
	} else {
		c.Stderr = os.Stderr
	}

	err := c.Run()
	if err != nil {
		if e := (&exec.ExitError{}); errors.As(err, &e) {
			return stdoutBuf.Bytes(), e.ExitCode()
		}
		return stdoutBuf.Bytes(), 1
	}
	return stdoutBuf.Bytes(), 0
}
