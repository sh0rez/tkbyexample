package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/radovskyb/watcher"
)

func main() {
	render()

	w := watcher.New()
	defer w.Close()
	w.SetMaxEvents(1)
	w.AddRecursive("./src/examples")

	go func() {
		for {
			watch(w)
		}
	}()

	go gatsbyDevelop()

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func watch(w *watcher.Watcher) {
	select {
	case <-w.Event:
		render()
	case err := <-w.Error:
		log.Fatalln(err)
	case <-w.Closed:
		return
	}
}

func gatsbyDevelop() {
	cmd := exec.Command("yarn", "dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

func render() {
	examples := parseExamples()
	renderExamples(examples)
}

var docsPat = regexp.MustCompile("^\\s*(\\/\\/|#)\\s")
var dashPat = regexp.MustCompile("\\-+")

// Seg is a segment of an example
type Seg struct {
	Docs, DocsRendered              string
	Code, CodeRendered, CodeForJs   string
	CodeEmpty, CodeLeading, CodeRun bool
}

// Example is info extracted from an example file
type Example struct {
	ID, Name                    string
	GoCode, GoCodeHash, URLHash string
	Segs                        [][]*Seg
	PrevExample                 *Example
	NextExample                 *Example
}

func renderExamples(examples []*Example) {
	tmpl := template.Must(template.ParseFiles("src/example.tmpl"))

	for _, example := range examples {
		buf := bytes.Buffer{}
		tmpl.Execute(&buf, example)

		if err := ioutil.WriteFile(filepath.Join("dist", example.ID+".md"), buf.Bytes(), 0644); err != nil {
			log.Fatalln(err)
		}
	}
}

func parseExamples() []*Example {
	var exampleNames []string
	for _, line := range readLines("examples.txt") {
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
		sourcePaths := mustGlob("src/examples/" + exampleID + "/*")

		for _, sourcePath := range sourcePaths {
			sourceSegs := parseAndRenderSegs(sourcePath)
			example.Segs = append(example.Segs, sourceSegs)
		}

		examples = append(examples, &example)
	}

	return examples
}

func parseSegs(sourcePath string) []*Seg {
	var (
		lines  []string
		source []string
	)
	// Convert tabs to spaces for uniform rendering.
	for _, line := range readLines(sourcePath) {
		lines = append(lines, strings.Replace(line, "\t", "    ", -1))
		source = append(source, line)
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
	return segs
}

func parseAndRenderSegs(sourcePath string) []*Seg {
	segs := parseSegs(sourcePath)

	for _, seg := range segs {
		if seg.Docs != "" {
			seg.DocsRendered = seg.Docs
		}
		if seg.Code != "" {
			// seg.CodeRendered = fmt.Sprintf("<pre>%s</pre>", seg.Code)
			seg.CodeRendered = seg.Code
		}
	}

	return segs
}

func readLines(path string) []string {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	return strings.Split(string(src), "\n")
}

func mustGlob(glob string) []string {
	paths, err := filepath.Glob(glob)
	if err != nil {
		log.Fatalln(err)
	}
	return paths
}
