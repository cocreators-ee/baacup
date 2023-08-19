//go:build windows

package main

import (
	"golang.org/x/sys/windows"
)

func getPlatform() string {
	return "windows"
}

func getBasePath() string {
	path, err := windows.KnownFolderPath(windows.FOLDERID_Documents, 0)
	if err != nil {
		panic(err)
	}
	return path
}
