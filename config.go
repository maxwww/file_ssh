package main

import (
	"os"
	"strconv"
)

type Config struct {
	Host                    string
	Username                string
	Port                    string
	KeyPath                 string
	AllFileNamesFilePath    string
	FilesPerSession         int
	RemoteServerFolder      string
	RemoteServerFilesFolder string
	LogFilePath             string
	HomeFilesFolder         string
}

func NewConfig() (*Config, error) {
	filesPerSession, err := strconv.Atoi(os.Getenv("FILES_PER_SESSION"))
	if err != nil {
		return nil, err
	}
	return &Config{
		Host:                    os.Getenv("HOST"),
		Username:                os.Getenv("REMOTE_USER"),
		Port:                    os.Getenv("PORT"),
		KeyPath:                 os.Getenv("PRIVATE_KEY_PATH"),
		AllFileNamesFilePath:    os.Getenv("ALL_FILE_NAMES_PATH"),
		FilesPerSession:         filesPerSession,
		RemoteServerFolder:      os.Getenv("REMOTE_SERVER_FOLDER"),
		RemoteServerFilesFolder: os.Getenv("REMOTE_SERVER_FILES_FOLDER"),
		LogFilePath:             os.Getenv("LOG_FILE_PATH"),
		HomeFilesFolder:         os.Getenv("HOME_FILES_FOLDER"),
	}, nil
}
