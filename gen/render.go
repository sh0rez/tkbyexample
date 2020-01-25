package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

func render() error {
	if err := os.MkdirAll("dist", os.ModePerm); err != nil {
		return err
	}

	examples, err := loadExamples()
	if err != nil {
		return err
	}

	return renderExamples(examples)
}

var docsPat = regexp.MustCompile("^\\s*(\\/\\/|#)\\s")
var dashPat = regexp.MustCompile("\\-+")

// Seg is a segment of an example
type Seg struct {
	Docs                            string
	Code                            string
	CodeEmpty, CodeLeading, CodeRun bool
}

// Example is info extracted from an example file
type Example struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Ignore      []string `yaml:"ignore"`

	ID   string   `yaml:"-"`
	Segs [][]*Seg `yaml:"-"`
}

func (e Example) CodeRaw() string {
	buf := ""
	for _, seg := range e.Segs {
		buf += joinSeg(seg) + "\n"
	}
	buf = strings.TrimSuffix(buf, "\n")

	return base64.StdEncoding.EncodeToString([]byte(buf))
}

func joinSeg(seg []*Seg) string {
	buf := ""
	for _, s := range seg {
		if s.CodeEmpty {
			continue
		}
		buf += s.Code + "\n"
	}
	return strings.TrimSuffix(buf, "\n")
}

func renderExamples(examples []*Example) error {
	tmpl, err := template.ParseFiles("src/example.tmpl")
	if err != nil {
		return err
	}

	for _, example := range examples {
		buf := bytes.Buffer{}
		tmpl.Execute(&buf, example)

		if err := ioutil.WriteFile(filepath.Join("dist", example.ID+".md"), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", example.ID, err)
		}
	}
	return nil
}

func loadExamples() ([]*Example, error) {
	var dirs []string
	root, err := filepath.Abs("./src/examples")
	if err != nil {
		return nil, err
	}

	filepath.Walk("src/examples", func(path string, fi os.FileInfo, err error) error {
		path, _ = filepath.Abs(path)
		if filepath.Base(path) != "x.yml" {
			return nil
		}

		dirs = append(dirs, filepath.Dir(path))
		return nil
	})

	examples := make([]*Example, 0, len(dirs))
	for _, d := range dirs {
		data, err := ioutil.ReadFile(filepath.Join(d, "x.yml"))
		if err != nil {
			return nil, err
		}

		rel, err := filepath.Rel(root, d)
		if err != nil {
			return nil, err
		}

		e := Example{
			ID: rel,
		}
		if err := yaml.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("parsing config for %s: %w", d, err)
		}

		e.Segs = make([][]*Seg, 0)
		sourcePaths, err := filepath.Glob(d + "/*")
		if err != nil {
			return nil, err
		}

		for _, sourcePath := range sourcePaths {
			if filepath.Base(sourcePath) == "x.yml" {
				continue
			}
			sourceSegs, err := parseSegs(sourcePath)
			if err != nil {
				return nil, err
			}
			e.Segs = append(e.Segs, sourceSegs)
		}

		examples = append(examples, &e)
	}

	return examples, nil
}

func parseSegs(sourcePath string) ([]*Seg, error) {

	var lines []string
	lines, err := readLines(sourcePath)
	if err != nil {
		return nil, err
	}

	// Convert tabs to spaces for uniform rendering.
	for i, line := range lines {
		lines[i] = strings.Replace(line, "\t", "    ", -1)
	}
	segs := []*Seg{}
	lastSeen := ""
	for _, line := range lines {
		if line == "" {
			lastSeen = ""
			continue
		}
		matchDocs := docsPat.MatchString(line)
		matchCode := !matchDocs
		newDocs := (lastSeen == "") || ((lastSeen != "docs") && (segs[len(segs)-1].Docs != ""))
		newCode := (lastSeen == "") || ((lastSeen != "code") && (segs[len(segs)-1].Code != ""))
		if matchDocs {
			trimmed := docsPat.ReplaceAllString(line, "")
			if newDocs {
				newSeg := Seg{Docs: trimmed, Code: ""}
				segs = append(segs, &newSeg)
			} else {
				segs[len(segs)-1].Docs = segs[len(segs)-1].Docs + "\n" + trimmed
			}
			lastSeen = "docs"
		} else if matchCode {
			if newCode {
				newSeg := Seg{Docs: "", Code: line}
				segs = append(segs, &newSeg)
			} else {
				segs[len(segs)-1].Code = segs[len(segs)-1].Code + "\n" + line
			}
			lastSeen = "code"
		}
	}
	for i, seg := range segs {
		seg.CodeEmpty = (seg.Code == "")
		seg.CodeLeading = (i < (len(segs) - 1))
		seg.CodeRun = strings.Contains(seg.Code, "package main")
	}
	return segs, nil
}

func readLines(path string) ([]string, error) {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(src), "\n"), nil
}
