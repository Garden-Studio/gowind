//go:build windows

package main

import (
	"errors"
	"os"
	"os/exec"
)

// run starts the CLI with this process's streams and exits with its status.
func run(binary string, args []string) error {
	command := exec.Command(binary, args...)
	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := command.Run()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		os.Exit(exit.ExitCode())
	}
	return err
}
