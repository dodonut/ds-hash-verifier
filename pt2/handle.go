package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"
)

type File struct {
	File multipart.File
	Handler *multipart.FileHeader
}

func GetFileFromForm(r *http.Request) (*File, error) {
	//max value accepted for file 10 MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, err
	}

	file, handler, err := r.FormFile("file")

	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	return &File{
		File:    file,
		Handler: handler,
	}, nil
}

func SaveFile(file multipart.File) (*string, error) {
	now := time.Now().Format(time.RFC3339Nano)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}

	err := os.MkdirAll(now, 0755)
	if err != nil {
		fmt.Println("erro ao criar diretorio")
		return nil, err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.txt", now, now), buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("erro ao criar arquivo")
		return nil, err
	}

	fmt.Println("Sucesso ao gravar arquivo")
	return &now, nil
}

func readAsync(filepath, hash string) {
	lineChan := make(chan string, 1000)
	found := make(chan bool, 1)
	var wg sync.WaitGroup
	wg.Add(4)

	for i := 0; i < 4; i++ {
		go processLine(hash, lineChan, found, &wg)
	}
	go readFileAsync(filepath, lineChan)

	wg.Wait()
	close(found)
	ok := <-found
	if !ok {
		fmt.Println("Falha ao encontrar senha.")
	}

}

func processLine(hash string, lineChan chan string, found chan bool, wg *sync.WaitGroup) {
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
			}
		}
	}

}

func readFileAsync(filepath string, lineChan chan string) {
	scan, err := GetFileFromDisk(filepath)
	if err != nil {
		fmt.Println("Nao foi possivel abrir o arquivo.")
		return
	}

	for scan.Scan() {
		lineChan <- scan.Text()
	}

	close(lineChan)
}

func Process(hash string) {
	var wg sync.WaitGroup
	wg.Add(4)

	go readAsync()
}

