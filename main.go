package main

import (
	"bufio"
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
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dir := fmt.Sprintf("%s/denote/", home)
	notes, err := parse(dir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(graph(notes))
}

func parse(dir string) ([]note, error) {
	var notes []note
	idRe := regexp.MustCompile(`[0-9]+T[0-9]+`)
	titleRe := regexp.MustCompile(`--[\pL-]+`)
	tagsRe := regexp.MustCompile(`_[\pL]+`)
	linkRe := regexp.MustCompile(`denote:[0-9]+T[0-9]+`)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %q: %v", path, err)
		}

		id := idRe.FindString(d.Name())
		if id == "" {
			return nil
		}

		title := titleRe.FindString(d.Name())
		title = strings.TrimPrefix(title, "--")

		tags := tagsRe.FindAllString(d.Name(), -1)
		for i := 0; i < len(tags); i++ {
			tags[i] = strings.TrimPrefix(tags[i], "_")
		}

		f, err := os.Open(dir + d.Name())
		if err != nil {
			fmt.Printf("error reading file %q: %v", d.Name(), err)
			return nil
		}
		defer f.Close()

		var links []string
		s := bufio.NewScanner(f)
		for s.Scan() {
			matches := linkRe.FindAllString(s.Text(), -1)
			for i := 0; i < len(matches); i++ {
				matches[i] = strings.TrimPrefix(matches[i], "denote:")
			}
			links = append(links, matches...)
		}
		if err := s.Err(); err != nil {
			fmt.Printf("error reading file %q: %v", d.Name(), err)
			return nil
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
		b.WriteString(fmt.Sprintf(`"%s" [label="%s"];`, n.id, n.title))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf(`"%s" -> {`, n.id))
		for _, l := range n.links {
			b.WriteString(fmt.Sprintf(`"%s" `, l))
		}
		b.WriteString("}\n")
	}
	b.WriteString("}")
	return b.String()
}
