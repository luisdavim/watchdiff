package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/luisdavim/watchdiff/pkg/watcher"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	opts := &watcher.Options{}
	var (
		colorOpt string
		execMode bool
	)

	var rootCmd = &cobra.Command{
		Use:   "watchdiff [command]",
		Short: "Watch a command and print colorized diffs of the output",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			switch colorOpt {
			case "always":
				opts.ColorEnabled = true
			case "never":
				opts.ColorEnabled = false
			case "auto":
				if info, err := os.Stdout.Stat(); err == nil {
					if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
						opts.ColorEnabled = true
					}
				}
			default:
				return fmt.Errorf("invalid value for color flag: %q", colorOpt)
			}

			if execMode {
				opts.Shell = ""
			}

			return watcher.Run(cmd.Context(), opts, args)
		},
	}

	rootCmd.Flags().DurationVarP(&opts.Interval, "interval", "n", 2*time.Second, "Interval between updates (e.g. 2s, 500ms)")
	rootCmd.Flags().IntVarP(&opts.ContextLines, "context", "c", 4, "Number of context lines for diff")
	rootCmd.Flags().BoolVarP(&opts.Quiet, "quiet", "q", false, "Suppress heartbeat dots")
	rootCmd.Flags().BoolVarP(&opts.IncludeStderr, "include-stderr", "e", true, "Include stderr in the diff comparison")
	rootCmd.Flags().StringVarP(&colorOpt, "color", "C", "auto", "Print colorized diffs (valid values are: auto, always or never)")
	rootCmd.Flags().BoolVarP(&execMode, "exec", "x", false, "Run the command directly, not through the shell")
	rootCmd.Flags().StringVarP(&opts.Shell, "shell", "s", "sh", "Specify the shell to use")

	return rootCmd
}
