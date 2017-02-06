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

const (
	htmlFlags = 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	markdownExtensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS
)

type Document struct {
	Title         string `fm:"title"`
	Content       string `fm:"content"`
	MarkdownTheme string
	TypekitKitID  string
	CodeTheme     string
}

type MarkdownHandlerOptions struct {
	DocRoot       string
	DocExtension  string
	DirIndex      string
	MarkdownTheme string
	TypekitKitID  string
	CodeTheme     string
}

func markdownHandleFunc(opts MarkdownHandlerOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path, err := filepath.Abs(filepath.Join(opts.DocRoot, r.URL.Path))
		if err != nil {
			serveInternalError(w, r, opts)
			fmt.Fprintf(os.Stderr, "Error finding absolute path (%s)", err)
			return
		}

		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			servePageNotFound(w, r, opts)
			return
		}

		if info.IsDir() {
			serveDirectory(w, r, path, opts)
			return
		}

		// Serve file
		serveFile(w, r, path, opts)
	}
}

func serveDirectory(w http.ResponseWriter, r *http.Request, path string, opts MarkdownHandlerOptions) {
	indexFile := filepath.Join(path, opts.DirIndex+opts.DocExtension)
	serveFile(w, r, indexFile, opts)
}

func serveFile(w http.ResponseWriter, r *http.Request, path string, opts MarkdownHandlerOptions) {
	ext := filepath.Ext(path)

	if ext == opts.DocExtension {
		serveMarkdown(w, r, path, opts)
	} else {
		http.ServeFile(w, r, path)
	}
}

func serveMarkdown(w http.ResponseWriter, r *http.Request, path string, opts MarkdownHandlerOptions) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file (%s).\n", err)
		serveInternalError(w, r, opts)
		return
	}

	doc := Document{
		MarkdownTheme: opts.MarkdownTheme,
		TypekitKitID:  opts.TypekitKitID,
		CodeTheme:     opts.CodeTheme,
	}

	err = frontmatter.Unmarshal(data, &doc)
	if err != nil {
		serveInternalError(w, r, opts)
		fmt.Fprintf(os.Stderr, "Error unmarshalling frontmatter (%s).\n", err)
	}

	doc.Content = parseMarkdown(doc.Content)
	serveDocument(w, r, doc)
}

func parseMarkdown(input string) string {
	// set up the HTML renderer
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")
	outputBytes := blackfriday.MarkdownOptions([]byte(input), renderer, blackfriday.Options{
		Extensions: markdownExtensions})
	return string(outputBytes)
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

func servePageNotFound(w http.ResponseWriter, r *http.Request, opts MarkdownHandlerOptions) {
	w.WriteHeader(http.StatusNotFound)
	doc := Document{
		MarkdownTheme: opts.MarkdownTheme,
		CodeTheme:     opts.CodeTheme,
		Title:         "Page Not Found",
		Content:       "<h1>Page Not Found</h1>",
	}
	serveDocument(w, r, doc)
}

func serveInternalError(w http.ResponseWriter, r *http.Request, opts MarkdownHandlerOptions) {
	w.WriteHeader(http.StatusInternalServerError)
	doc := Document{
		MarkdownTheme: opts.MarkdownTheme,
		CodeTheme:     opts.CodeTheme,
		Title:         "Invalid Request",
		Content:       "<h1>Invalid Request</h1>",
	}
	serveDocument(w, r, doc)
}
