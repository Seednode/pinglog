/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func removeExisting(logFile string) error {
	if ForceOverwrite {
		_, err := fmt.Print("File " + logFile + " already exists. Remove? (y/N) ")
		if err != nil {
			return err
		}

		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		if input.Text() != "y" {
			os.Exit(1)
		}
	}

	err := os.Remove(logFile)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	return nil
}

func checkExisting(logFile string) error {
	_, err := os.Stat(logFile)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return nil
	case err != nil:
		return err
	}

	err = removeExisting(logFile)
	if err != nil {
		return err
	}

	return nil
}

func teeOutput(logFilePtr *os.File) (func(), error) {
	stdOut := os.Stdout

	multiWriter := io.MultiWriter(stdOut, logFilePtr)

	reader, writer, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	os.Stdout = writer

	exit := make(chan bool)

	go func() {
		_, err = io.Copy(multiWriter, reader)
		if err != nil {
			log.Fatal(err)
		}

		exit <- true
	}()

	return func() {
		_ = writer.Close()
		<-exit
		_ = logFilePtr.Close()
	}, nil
}

func logOutput(logFile string) (func(), error) {
	if strings.HasPrefix(logFile, "~/") {
		dirname, _ := os.UserHomeDir()
		logFile = filepath.Join(dirname, logFile[2:])
	}

	err := checkExisting(logFile)
	if err != nil {
		return nil, err
	}

	logFilePtr, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}

	function, err := teeOutput(logFilePtr)
	if err != nil {
		return nil, err
	}

	return function, nil
}
