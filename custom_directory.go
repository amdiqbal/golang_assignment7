package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

type FileSystem struct {
	fs http.FileSystem
}

func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			return nil, err
		}
	}
	return f, nil
}

func main() {
	port := flag.String("p", "999", "port to serve on")
	directory := flag.String("d", ".", "serve a directory")
	flag.Parse()

	fileServer := http.FileServer(FileSystem{http.Dir(*directory)})
	http.Handle("C:Go/", http.StripPrefix(strings.TrimRight("/assignment7/", "/"), fileServer))

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
