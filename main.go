package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type note struct {
	id    string
	title string
	tags  []string
	links []string
}

func main() {
	var dir = "/home/davide/denote/"
	notes, err := parse(dir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(notes)
}

func parse(dir string) ([]note, error) {
	var notes []note
	idre := regexp.MustCompile(`[0-9]+T[0-9]+`)
	titlere := regexp.MustCompile(`--[\p{L}-]+`)
	tagsre := regexp.MustCompile(`_[\p{L}]+`)
	linkre := regexp.MustCompile(`denote:[0-9]+T[0-9]+`)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Failure accessing a path %q: %v\n", path, err)
			return err
		}

		fname := d.Name()

		id := idre.FindString(fname)
		if id == "" {
			return nil
		}

		title := titlere.FindString(fname)
		title = strings.TrimPrefix(title, "--")

		tags := tagsre.FindAllString(fname, -1)
		for i := 0; i < len(tags); i++ {
			tags[i] = strings.TrimPrefix(tags[i], "_")
		}

		dat, err := os.ReadFile(dir + fname)
		if err != nil {
			return nil
		}
		links := linkre.FindAllString(string(dat), -1)
		for i := 0; i < len(links); i++ {
			links[i] = strings.TrimPrefix(links[i], "denote:")
		}

		notes = append(notes, note{id, title, tags, links})

		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return nil, err
	}

	return notes, nil
}
