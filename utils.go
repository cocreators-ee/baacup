//go:build !windows

package main

import (
	"os"
	"runtime"
)

func getPlatform() string {
	if runtime.GOOS == "darwin" {
		return "macos"
	}

	// With some luck this will give defaults that end up working on e.g. BSD as well
	return "linux"
}

func getBasePath() string {
	// Darwin, Linux, other *nix
	return os.ExpandEnv("$HOME/Baacup")
}
