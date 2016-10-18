package main

import (
	"log"
	"os"
)

func main() {
	l := log.New(os.Stderr, "", 0)
	if err := run(l); err != nil {
		log.Fatal(err)
	}
}

func run(l *log.Logger) error {
	return nil
}
