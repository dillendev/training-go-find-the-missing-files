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

func Search(root string, terms []string) chan (string) {
	matches := make(chan string)

	go func() {
		var wg sync.WaitGroup
		search(&wg, matches, root, terms)

		wg.Wait()

		close(matches)
	}()

	return matches
}

func search(wg *sync.WaitGroup, matches chan (string), root string, terms []string) {
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
				search(wg, matches, path, terms)
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
}

func findMatch(path string, terms []string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		var buf [4096]byte

		if _, err := reader.Read(buf[:]); err != nil {
			if errors.Is(err, io.EOF) {
				return "", errNoMatch
			}

			return "", err
		}

		for _, term := range terms {
			if bytes.Contains(buf[:], []byte(term)) {
				return path, nil
			}
		}
	}
}
