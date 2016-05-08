package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ericaro/frontmatter"
	"github.com/russross/blackfriday"
)

type Document struct {
	Title         string `fm:"title"`
	Content       string `fm:"content"`
	MarkdownTheme string
	CodeTheme     string
}

func markdownHandleFunc(c Config) httpHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path, err := filepath.Abs(filepath.Join(c.DocRoot, r.URL.Path))
		if err != nil {
			serveInternalError(w, r, c)
			fmt.Fprintf(os.Stderr, "Error finding absolute path (%s)", err)
			return
		}

		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			servePageNotFound(w, r, c)
			return
		}

		if info.IsDir() {
			serveDirectory(w, r, path, c)
			return
		}

		// Serve file
		serveFile(w, r, path, c)
	}
}

func serveDirectory(w http.ResponseWriter, r *http.Request, path string, c Config) {
	indexFile := filepath.Join(path, c.DirIndex+c.DocExtension)
	serveFile(w, r, indexFile, c)
}

func serveFile(w http.ResponseWriter, r *http.Request, path string, c Config) {
	ext := filepath.Ext(path)

	if ext == c.DocExtension {
		serveMarkdown(w, r, path, c)
	} else {
		http.ServeFile(w, r, path)
	}
}

func serveMarkdown(w http.ResponseWriter, r *http.Request, path string, c Config) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file (%s).\n", err)
		serveInternalError(w, r, c)
		return
	}

	doc := Document{
		MarkdownTheme: c.MarkdownTheme,
		CodeTheme:     c.CodeTheme,
	}

	err = frontmatter.Unmarshal(data, &doc)
	if err != nil {
		serveInternalError(w, r, c)
		fmt.Fprintf(os.Stderr, "Error unmarshalling frontmatter (%s).\n", err)
	}

	doc.Content = string(blackfriday.MarkdownCommon([]byte(doc.Content)))
	serveDocument(w, r, doc)
}

func serveDocument(w http.ResponseWriter, r *http.Request, doc Document) {
	templateData, err := Asset("assets/templates/document.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding template (%s).\n", err)
	}

	t, err := template.New("document").Parse(string(templateData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template (%s).\n", err)
	}

	t.Execute(w, doc)
}

func servePageNotFound(w http.ResponseWriter, r *http.Request, c Config) {
	w.WriteHeader(http.StatusNotFound)
	doc := Document{
		MarkdownTheme: c.MarkdownTheme,
		CodeTheme:     c.CodeTheme,
		Title:         "Page Not Found",
		Content:       "<h1>Page Not Found</h1>",
	}
	serveDocument(w, r, doc)
}

func serveInternalError(w http.ResponseWriter, r *http.Request, c Config) {
	w.WriteHeader(http.StatusInternalServerError)
	doc := Document{
		MarkdownTheme: c.MarkdownTheme,
		CodeTheme:     c.CodeTheme,
		Title:         "Invalid Request",
		Content:       "<h1>Invalid Request</h1>",
	}
	serveDocument(w, r, doc)
}
