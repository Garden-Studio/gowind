//go:build !windows

package main

import (
	"os"
	"syscall"
)

// run replaces this process with the CLI, so signals (a --watch stopped with
// Ctrl-C) and the exit status are its own.
func run(binary string, args []string) error {
	return syscall.Exec(binary, append([]string{binary}, args...), os.Environ())
}
