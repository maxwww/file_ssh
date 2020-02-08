package main

import (
	"os"
)

func logToFile(data []byte, config *Config) error {
	f, err := os.OpenFile(config.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	data = append(data, 10)
	if _, err := f.Write(data); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
