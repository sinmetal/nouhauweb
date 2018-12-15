package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/shurcooL/github_flavored_markdown"
)

func main() {
	fmt.Println(dirwalk("."))
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
		if strings.HasSuffix(file.Name(), ".md") {
			if err := markdownToHTML(filepath.Join(dir, file.Name())); err != nil {
				panic(err)
			}
		}
	}

	return paths
}

func markdownToHTML(filePath string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var bb bytes.Buffer
	bb.WriteString(`<html><head><meta charset="utf-8"><link href="/assets/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
	bb.Write(github_flavored_markdown.Markdown(content))
	bb.WriteString(`</article></body></html>`)

	if err := writeHTML(filePath, bb.Bytes()); err != nil {
		return err
	}

	return nil
}

func writeHTML(filePath string, data []byte) error {
	return ioutil.WriteFile(fmt.Sprintf("%s.html", filePath), data, 0666)
}
