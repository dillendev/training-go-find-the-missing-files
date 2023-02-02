package main

import (
	"flag"
	"fmt"

	"github.com/dillendev/training-go-find-the-missing-files/internal/grep"
)

func main() {
	flag.Parse()

	terms := flag.Args()
	matches := grep.Search("exampledata", terms)

	for match := range matches {
		fmt.Println(match.Path)
	}
}
