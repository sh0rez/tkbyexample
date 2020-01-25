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
)

func render() error {
	if err := os.MkdirAll("dist", os.ModePerm); err != nil {
		return err
	}

	examples, err := parseExamples()
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
	ID, Name string
	Segs     [][]*Seg
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

func parseExamples() ([]*Example, error) {
	var exampleNames []string

	lines, err := readLines("examples.txt")
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		if line != "" && !strings.HasPrefix(line, "#") {
			exampleNames = append(exampleNames, line)
		}
	}

	examples := make([]*Example, 0)
	for _, exampleName := range exampleNames {
		example := Example{Name: exampleName}

		exampleID := strings.ToLower(exampleName)
		exampleID = strings.Replace(exampleID, " ", "-", -1)
		exampleID = strings.Replace(exampleID, "/", "-", -1)
		exampleID = strings.Replace(exampleID, "'", "", -1)
		exampleID = dashPat.ReplaceAllString(exampleID, "-")
		example.ID = exampleID

		example.Segs = make([][]*Seg, 0)
		sourcePaths, err := filepath.Glob("src/examples/" + exampleID + "/*")
		if err != nil {
			return nil, err
		}

		for _, sourcePath := range sourcePaths {
			sourceSegs, err := parseSegs(sourcePath)
			if err != nil {
				return nil, err
			}
			example.Segs = append(example.Segs, sourceSegs)
		}

		examples = append(examples, &example)
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
	for _, line := range lines {
		lines = append(lines, strings.Replace(line, "\t", "    ", -1))
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
