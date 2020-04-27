package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func readSync(scan *bufio.Scanner, hash string) {
	for scan.Scan() {
		line := scan.Text()
		mdHash := md5.Sum([]byte(line))
		if hash == hex.EncodeToString(mdHash[:]) {
			fmt.Println("Encontrei: ", line)
			return
		}
	}
	fmt.Println("Falha ao encontrar senha")
}

func readAsync(f *os.File, scan *bufio.Scanner, hash string) {
	lineChan := make(chan string, 1000)
	found := make(chan bool, 1)
	var wg sync.WaitGroup
	wg.Add(4)

	for i := 0; i < 4; i++ {
		go processLine(f, hash, lineChan, found, &wg)
	}
	go readFileAsync(scan, lineChan)

	wg.Wait()
	close(found)
	ok := <-found
	if !ok {
		fmt.Println("Falha ao encontrar senha.")
	}

}

func processLine(f *os.File, hash string, lineChan chan string, found chan bool, wg *sync.WaitGroup) {
	for {
		select {
		case line, ok := <-lineChan:
			if !ok {
				wg.Done()
				return
			}
			mdHash := md5.Sum([]byte(line))
			if hex.EncodeToString(mdHash[:]) == hash {
				fmt.Println("Encontrei: ", line)
				found <- true
				f.Close()
			}
		}
	}

}

func readFileAsync(scan *bufio.Scanner, lineChan chan string) {
	for scan.Scan() {
		lineChan <- scan.Text()
	}
	close(lineChan)
}
func main() {
	file, err := os.Open("rockyou.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	hash := "fdaf6fee6b184c101e318114715b6814"
	start := time.Now()
	// readSync(scanner, hash)
	//b1tch3s fdaf6fee6b184c101e318114715b6814
	readAsync(file, scanner, hash)
	fmt.Printf("Tempo: %2f segundos", time.Since(start).Seconds())
}
