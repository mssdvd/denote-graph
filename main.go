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
	fmt.Println(graph(notes))
}

func parse(dir string) ([]note, error) {
	var notes []note
	idre := regexp.MustCompile(`[0-9]+T[0-9]+`)
	titlere := regexp.MustCompile(`--[\p{L}-]+`)
	tagsre := regexp.MustCompile(`_[\p{L}]+`)
	linkre := regexp.MustCompile(`denote:[0-9]+T[0-9]+`)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %q: %v", path, err)
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
			fmt.Printf("error reading file %q: %v", d.Name(), err)
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
		return nil, fmt.Errorf("error walking the path: %v", err)
	}

	return notes, nil
}

func graph(notes []note) string {
	var b strings.Builder

	b.WriteString("digraph denote {\n")

	for _, n := range notes {
		b.WriteString("\"")
		b.WriteString(n.id)
		b.WriteString("\" [label=\"")
		b.WriteString(n.title)
		b.WriteString("\"];\n")

		b.WriteString("\"")
		b.WriteString(n.id)
		b.WriteString("\" -> {")
		for _, l := range n.links {
			b.WriteString(" \"")
			b.WriteString(l)
			b.WriteString("\" ")
		}
		b.WriteString("}\n")
	}

	b.WriteString("}")
	return b.String()
}
