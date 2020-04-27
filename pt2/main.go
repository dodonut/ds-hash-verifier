package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	http.HandleFunc("/master/", func(writer http.ResponseWriter, request *http.Request) {
		//ctx, cancel := context.WithCancel(r.Context())

		hash := strings.TrimPrefix(request.URL.Path, "/master/")

		file, err := GetFileFromForm(request)
		if err != nil {
			fmt.Println("Nao foi possivel armazenar o arquivo.")
			return
		}

		fileName, err := SaveFile(file.File)
		if err != nil {
			fmt.Println("Nao foi possivel armazenar o arquivo.")
			return
		}

		err := partitionFile(*fileName)
		if err != nil {
			fmt.Println("Erro ao particionar arquivo")
		}

		Process(hash)
	})

	fmt.Println("Listening to 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erro ao abrir servidor")
	}


}