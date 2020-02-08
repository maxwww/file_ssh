package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/ssh"
	"math/rand"
)

func RemoteRun(command string, client *ssh.Client, config *Config) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var sessionStdout bytes.Buffer
	var sessionStderr bytes.Buffer
	session.Stdout = &sessionStdout
	session.Stderr = &sessionStderr

	err = session.Run(command)
	if err != nil {
		return err
	}

	if err := logToFile(sessionStderr.Bytes(), config); err != nil {
		return err

	}
	if err := logToFile(sessionStdout.Bytes(), config); err != nil {
		return err

	}

	return nil
}

func tarRemote(fileWithList string, client *ssh.Client, config *Config) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var sessionStdout bytes.Buffer
	var sessionStderr bytes.Buffer
	session.Stdout = &sessionStdout
	session.Stderr = &sessionStderr

	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	archiveFileName := hex.EncodeToString(randBytes) + ".tar"
	archivePath := fmt.Sprintf("%s/%s", config.RemoteServerFolder, archiveFileName)
	command := fmt.Sprintf("cd %s && tar --ignore-failed-read -cf %s -T %s", config.RemoteServerFilesFolder, archivePath, fileWithList)

	if err = RemoteRun(command, client, config); err != nil {
		return "", err
	}
	return archivePath, nil
}

func RemoteRemoveFile(file string, client *ssh.Client, config *Config) error {
	command := fmt.Sprintf("rm -rf %s", file)
	return RemoteRun(command, client, config)
}
