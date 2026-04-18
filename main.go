package main

import (
	"os"

	"github.com/luisdavim/watchdiff/cmd"
)

func main() {
	rootCmd := cmd.New()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
