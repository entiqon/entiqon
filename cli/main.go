package main

import (
	"fmt"
	"os"
	"os/exec"
)

var scripts = map[string]string{
	"gsu":     "bin/gsu.sh",
	"gotestx": "bin/gotestx.sh",
	"gcr":     "bin/gcr.sh",
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: entiqon-cli <command> [args]")
		os.Exit(1)
	}

	cmd, ok := scripts[os.Args[1]]
	if !ok {
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}

	c := exec.Command(cmd, os.Args[2:]...)
	c.Stdout, c.Stderr, c.Stdin = os.Stdout, os.Stderr, os.Stdin
	if err := c.Run(); err != nil {
		_, err = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
