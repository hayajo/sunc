// +build linux

package main

import (
	"fmt"
	"os"
	"runtime"
)

func init() {
	if os.Args[0] == INIT_COMMAND {

		runtime.GOMAXPROCS(1)
		runtime.LockOSThread()

		proc := Process(os.Args[1], os.Args[2:]...)
		if err := proc.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Fork error: %v\n", err)
			os.Exit(1)
		}

		panic("this line should have never been reached")
	}
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		args = []string{os.Getenv("SHELL")}
	}

	ct := Container(args[0], args[1:]...)
	if err := ct.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start container: %v\n", err)
		os.Exit(1)
	}
}
