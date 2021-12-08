package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

type SpaHandler struct {
	staticPath string
	indexPath  string
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("No storage location specified.")
		return
	}

	dropLocation := os.Args[1]

	router := mux.NewRouter()

	router.HandleFunc("/api/health-check", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "of great health!")
	})

	router.HandleFunc("/api/deposit", deposit(dropLocation))

	application := SpaHandler{staticPath: "static", indexPath: "index.html"}
	router.PathPrefix("/").Handler(application)

	server := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Listening on port 8080...")
	fmt.Printf("Clients can add files at [[ %s:8080/app/deposit ]]\n", getLocalIP())
	log.Fatal(server.ListenAndServe())
}

func deposit(location string) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
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

			fmt.Printf("Uploaded %s from %s\n", handler.Filename, request.RemoteAddr)
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

func (h SpaHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	path, err := filepath.Abs(request.URL.Path)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	servePath := filepath.Join(h.staticPath, path)

	_, err = os.Stat(servePath)
	if os.IsNotExist(err) {
		// serve index file, handing off route to react SPA
		http.ServeFile(response, request, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(response, request)
}
