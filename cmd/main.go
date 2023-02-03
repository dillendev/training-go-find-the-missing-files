package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dillendev/training-go-find-the-missing-files/internal/grep"
)

func main() {
	flag.Parse()

	terms := flag.Args()

	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	matches := grep.Search(root, terms)

	for path := range matches {
		fmt.Println(path)
	}
}
