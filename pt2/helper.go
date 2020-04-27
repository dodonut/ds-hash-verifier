package main

import (
	"bufio"
	"os"
	"os/exec"
)

func GetFileFromDisk(fileName string) (*bufio.Scanner, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	return bufio.NewScanner(file), nil
}

func partitionFile(filename string) error {
	return exec.Command("sh", "partition.sh", filename).Run()
}