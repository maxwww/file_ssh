package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func copyToServer(input io.Reader, config *Config) (string, error) {
	output, err := ioutil.TempFile("tmp", "")
	if err != nil {
		return "", err
	}
	defer os.Remove(output.Name())

	if _, err := io.Copy(output, input); err != nil {
		return "", err
	}

	fileInfo, err := output.Stat()
	if err != nil {
		return "", err
	}

	fullRemoteFilePath := fmt.Sprintf("%s/%s", config.RemoteServerFolder, fileInfo.Name())
	dstPath := fmt.Sprintf("%s@%s:%s", config.Username, config.Host, fullRemoteFilePath)

	cmd := exec.Command("scp", "-i", config.KeyPath, output.Name(), dstPath)
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	if err := output.Close(); err != nil {
		return "", err
	}

	return fullRemoteFilePath, nil
}

func copyFromServer(remoteFilePath string, config *Config) error {
	src := fmt.Sprintf("%s@%s:%s", config.Username, config.Host, remoteFilePath)

	cmd := exec.Command("scp", "-i", config.KeyPath, src, config.HomeFilesFolder)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
