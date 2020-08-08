package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("No storage location specified.")
		return
	}

	dropLocation := os.Args[1]

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	http.HandleFunc("/health-check", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "of great health!")
	})

	http.HandleFunc("/deposit", deposit(dropLocation))

	fmt.Println("Listening on port 8080...")
	fmt.Printf("Clients can add files at [[ %s:8080/deposit ]]\n", getLocalIP())
	http.ListenAndServe(":8080", nil)
}

func deposit(location string) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			http.ServeFile(response, request, "./static/deposit.html")
		} else if request.Method == "POST" {
			request.ParseMultipartForm(32 << 20)
			uploadFile, handler, err := request.FormFile("uploadfile")
			if err != nil {
				fmt.Println(err)
				http.Error(response, "The server doesn't like this file", http.StatusBadRequest)
				return
			}

			defer uploadFile.Close()

			fmt.Fprintf(response, "%s uploaded!", handler.Filename)

			localFile, err := os.OpenFile(location+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0664)
			if err != nil {
				fmt.Println(err)
				return
			}

			defer localFile.Close()
			io.Copy(localFile, uploadFile)

			fmt.Printf("Uploaded %s from %s", handler.Filename, request.RemoteAddr)
		} else {
			http.Error(response, "Is not how this works.", http.StatusNotFound)
			return
		}
	}
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "208.67.222.222:80")
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
