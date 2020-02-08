package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func ExtractFile(archivePath string, config *Config, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("starting extracting %s file on local host\n", archivePath)

	fmt.Println(archivePath)

	s := strings.Split(archivePath, "/")

	localArchivePath := fmt.Sprintf("%s/%s", config.HomeFilesFolder, s[len(s)-1])
	command := []string{
		"-xf",
		localArchivePath,
		"-C",
		config.HomeFilesFolder,
	}

	cmd := exec.Command("tar", command...)
	if err := cmd.Run(); err != nil {
		_ = logToFile([]byte(err.Error()), config)
	} else {
		fmt.Printf("successfully extracted %s file on local host\n", archivePath)
	}

	cmd2 := exec.Command("rm", "-rf", localArchivePath)
	if err := cmd2.Run(); err != nil {
		_ = logToFile([]byte(err.Error()), config)
	}
}
