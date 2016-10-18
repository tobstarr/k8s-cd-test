package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	l := log.New(os.Stderr, "", 0)
	if err := run(l); err != nil {
		log.Fatal(err)
	}
}

func run(l *log.Logger) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	if !strings.Contains(port, ":") {
		port = "0.0.0.0:" + port
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/clicked", clicked)
	l.Printf("running on port %q", port)
	return http.ListenAndServe(port, mux)
}

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `<html><body><a href="#" onclick="javascript:document.location = '/clicked';">Click me</a>`)
}

func clicked(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `<html><body><div id="clicked"><h2>you clicked!</h2></div>`)
}
