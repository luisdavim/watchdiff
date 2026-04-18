package watcher

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aymanbagabas/go-udiff"
)

const (
	timeFmt = "15:04:05"
)

type Options struct {
	Interval      time.Duration
	ContextLines  int
	Quiet         bool
	IncludeStderr bool
	ColorEnabled  bool
	Shell         string
}

func Run(ctx context.Context, opts *Options, args []string) error {
	var clrBlue, clrYellow, clrGreen, clrGray, clrRed, clrCyan, clrReset string
	if opts.ColorEnabled {
		clrBlue = "\033[1;34m"
		clrYellow = "\033[1;33m"
		clrGreen = "\033[1;32m"
		clrGray = "\033[90m"
		clrRed = "\033[1;31m"
		clrCyan = "\033[36m"
		clrReset = "\033[0m"
	}

	fmt.Printf("%sMonitoring:%s %s\n", clrGreen, clrReset, strings.Join(args, " "))
	fmt.Printf("%sInterval: %s | Context: %d | Stderr: %v%s\n", clrGray, opts.Interval, opts.ContextLines, opts.IncludeStderr, clrReset)
	fmt.Println(clrBlue + strings.Repeat("-", 3) + clrReset)

	if opts.Shell == "" {
		opts.Shell = args[0]
		args = args[1:]
	} else {
		args = []string{"-c", strings.Join(args, " ")}
	}

	lastOutput, lastExit := execute(ctx, opts.Shell, args, opts.IncludeStderr)
	fmt.Printf("%sINITIAL OUTPUT (%s):%s\n%s\n", clrBlue, time.Now().Format(timeFmt), clrReset, string(lastOutput))
	fmt.Println(clrBlue + strings.Repeat("-", 3) + clrReset)

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\n%sShutting down...%s\n", clrYellow, clrReset)
			return nil
		case <-ticker.C:
			currentOutput, currentExit := execute(ctx, opts.Shell, args, opts.IncludeStderr)
			if ctx.Err() != nil {
				return nil
			}

			edits := udiff.Bytes(lastOutput, currentOutput)
			if len(edits) > 0 || currentExit != lastExit {
				fmt.Printf("\n%sCHANGE DETECTED @ %s%s\n", clrYellow, time.Now().Format(timeFmt), clrReset)

				if currentExit != lastExit {
					clr := clrRed
					if currentExit == 0 {
						clr = clrGreen
					}
					fmt.Printf("%sExit Code Changed: %d -> %d%s\n", clr, lastExit, currentExit, clrReset)
				}

				// Generate diff
				diffText, err := udiff.ToUnified("Old", "New", string(lastOutput), edits, opts.ContextLines)
				if err != nil {
					return fmt.Errorf("faild to compute diff: %w", err)
				}

				if diffText != "" {
					printColorizedDiff(diffText, clrRed, clrGreen, clrCyan, clrGray, clrReset)
				}

				lastOutput = currentOutput
				lastExit = currentExit
			} else if !opts.Quiet {
				fmt.Printf("%s.%s", clrGray, clrReset)
			}
		}
	}
}

func execute(ctx context.Context, cmd string, args []string, captureStderr bool) ([]byte, int) {
	var stdoutBuf bytes.Buffer
	c := exec.CommandContext(ctx, cmd, args...)
	// c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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

func printColorizedDiff(d, removed, added, header, text, reset string) {
	for line := range strings.Lines(d) {
		if line == "" || line == "--- Old\n" || line == "+++ New\n" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "+"):
			fmt.Printf("%s%s%s", added, line, reset)
		case strings.HasPrefix(line, "-"):
			fmt.Printf("%s%s%s", removed, line, reset)
		case strings.HasPrefix(line, "@@"):
			fmt.Printf("%s%s%s", header, line, reset)
		default:
			fmt.Printf("%s%s%s", text, line, reset)
		}
	}
}
