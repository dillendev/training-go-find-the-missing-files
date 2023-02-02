package grep

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var errNoMatch = errors.New("no matches")

func Search(root string, terms []string) chan (Match) {
	matches := make(chan Match)

	go func() {
		search(root, terms, matches)
		close(matches)
	}()

	return matches
}

func search(root string, terms []string, matches chan (Match)) {
	var wg sync.WaitGroup

	entries, err := os.ReadDir(root)
	if err != nil {
		log.Printf("error reading directory %s: %s", root, err.Error())
		return
	}

	for _, entry := range entries {
		path := filepath.Join(root, entry.Name())

		if entry.IsDir() {
			wg.Add(1)

			go func() {
				defer wg.Done()
				search(path, terms, matches)
			}()
			continue
		}

		if !entry.Type().IsRegular() {
			continue
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			match, err := findMatch(path, terms)

			if err != nil {
				if errors.Is(err, errNoMatch) {
					return
				}

				log.Printf("error searching for match %s: %s", path, err.Error())
				return
			}

			matches <- match
		}()
	}

	wg.Wait()
}

func findMatch(path string, terms []string) (Match, error) {
	file, err := os.Open(path)
	if err != nil {
		return Match{}, err
	}

	defer file.Close()

	var buff [2048]byte
	reader := bufio.NewReader(file)

	if _, err := reader.Read(buff[:]); err != nil {
		if errors.Is(err, io.EOF) {
			return Match{}, errNoMatch
		}

		return Match{}, err
	}

	for _, term := range terms {
		if bytes.Contains(buff[:], []byte(term)) {
			return Match{Path: path}, nil
		}
	}

	return Match{}, errNoMatch
}