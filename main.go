package main

import (
	"devkit/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
