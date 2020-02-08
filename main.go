package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"sync"
)

var fileHandled = 0

func main() {
	err := godotenv.Load()
	var wg sync.WaitGroup
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config, err := NewConfig()
	if err != nil {
		log.Fatal("Error to getting config")
	}

	allFileNamesFile, err := os.Open(config.AllFileNamesFilePath)
	if err != nil {
		log.Fatal("error to open file: ", err)
	}
	defer allFileNamesFile.Close()

	client, err := getClient(config)
	if err != nil {
		log.Fatal("error to get client: ", err)
	}
	defer client.Close()

	totalLines, err := lineCounter(allFileNamesFile)
	if err != nil {
		log.Fatal("error to count lines: ", err)
	}
	allFileNamesFile.Seek(0, io.SeekStart)

	scanner := bufio.NewScanner(allFileNamesFile)
	var buffer bytes.Buffer
	countLines := 0
	for scanner.Scan() {
		countLines++
		buffer.Write(scanner.Bytes())
		buffer.Write([]byte{10})
		if countLines%config.FilesPerSession == 0 || countLines == totalLines {
			err := handleFiles(&buffer, config, client, &wg)
			if err != nil {
				log.Fatal("error to handle files: ", err)
			}
		}
	}

	wg.Wait()
}

func handleFiles(input io.Reader, config *Config, client *ssh.Client, wg *sync.WaitGroup) error {
	fileHandled++
	fmt.Printf("strting handling #%d work", fileHandled)

	fmt.Println("starting copping list of files to remote server")
	fileWithList, err := copyToServer(input, config)
	if err != nil {
		return err
	}
	fmt.Printf("successfully copped %s file to remote server\n", fileWithList)

	fmt.Println("starting tarring list of files at remote server")
	archivePath, err := tarRemote(fileWithList, client, config)
	if err != nil {
		return err
	}
	fmt.Printf("successfully tarred %s file at remote server\n", archivePath)

	fmt.Printf("starting downloading archive %s from remote server\n", archivePath)
	if err = copyFromServer(archivePath, config); err != nil {
		return err
	}
	fmt.Printf("successfully downloaded %s archive from remote server\n", archivePath)

	fmt.Printf("starting removing %s file from remote server\n", fileWithList)
	if err = RemoteRemoveFile(fileWithList, client, config); err != nil {
		return err
	}
	fmt.Printf("successfully removed %s file from remote server\n", fileWithList)

	fmt.Printf("starting removing %s file from remote server\n", archivePath)
	if err = RemoteRemoveFile(archivePath, client, config); err != nil {
		return err
	}
	fmt.Printf("successfully removed %s file from remote server\n", archivePath)

	wg.Add(1)
	go ExtractFile(archivePath, config, wg)

	fmt.Printf("successfully handled #%d work\n", fileHandled)
	return nil
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
