package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func logOutput(logFile string) func() {
	_, err := os.Stat(logFile)
	if errors.Is(err, os.ErrNotExist) == false {
		fmt.Print("File " + logFile + " already exists. Remove? (y/N) ")

		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		if input.Text() != "y" {
			os.Exit(1)
		}

		err := os.Remove(logFile)
		if err != nil {
			panic(err)
		}
	}

	logFilePtr, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}

	out := os.Stdout

	multiWriter := io.MultiWriter(out, logFilePtr)

	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	os.Stdout = writer
	os.Stderr = writer

	exit := make(chan bool)
	go func() {
		_, _ = io.Copy(multiWriter, reader)
		exit <- true
	}()

	return func() {
		_ = writer.Close()
		<-exit
		_ = logFilePtr.Close()
	}
}
