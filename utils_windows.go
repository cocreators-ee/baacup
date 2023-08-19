//go:build windows

package main

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func getPlatform() string {
	return "windows"
}

func getBasePath() string {
	myDocuments, err := windows.KnownFolderPath(windows.FOLDERID_Documents, 0)
	if err != nil {
		panic(err)
	}

	// Save this so we can reference it in rules
	os.Setenv("MyDocuments", myDocuments)

	return filepath.Join(myDocuments, "Baacup")
}
