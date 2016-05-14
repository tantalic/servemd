package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jawher/mow.cli"
)

const (
	Version = "0.3.2"
)

type httpHandleFunc func(w http.ResponseWriter, r *http.Request)

func main() {
	app := cli.App("servemd", "a simple http server for markdown content")
	app.Version("v version", Version)
	app.Spec = "[OPTIONS] [DIR]"

	var (
		// HTTP Options
		host = app.String(cli.StringOpt{
			Name:   "a host",
			Desc:   "Host/IP address to listen on for HTTP",
			Value:  "",
			EnvVar: "HOST",
		})
		port = app.Int(cli.IntOpt{
			Name:   "p port",
			Desc:   "TCP PORT to listen on for HTTP",
			Value:  3000,
			EnvVar: "PORT",
		})
		users = app.Strings(cli.StringsOpt{
			Name:   "u auth",
			Desc:   "Username and password for basic authentication in the form of user:pass",
			EnvVar: "BASIC_AUTH",
		})

		// Content
		dir = app.String(cli.StringArg{
			Name:   "DIR",
			Desc:   "Directory to serve content from",
			Value:  ".",
			EnvVar: "DOCUMENT_ROOT",
		})
		extension = app.String(cli.StringOpt{
			Name:   "e extension",
			Desc:   "Extension used for markdown files",
			Value:  ".md",
			EnvVar: "DOCUMENT_EXTENSION",
		})
		index = app.String(cli.StringOpt{
			Name:   "i index",
			Desc:   "Filename (without extension) to use for directory index",
			Value:  "index",
			EnvVar: "DIRECTORY_INDEX",
		})

		// Theme
		markdownTheme = app.String(cli.StringOpt{
			Name:   "m markdown-theme",
			Desc:   "Theme to use for styling markdown html",
			Value:  "clean",
			EnvVar: "MARKDOWN_THEME",
		})
		codeTheme = app.String(cli.StringOpt{
			Name:   "c code-theme",
			Desc:   "Highlight.js theme to use for syntax highlighting",
			Value:  "",
			EnvVar: "CODE_THEME",
		})
	)

	app.Action = func() {
		// Static Asset Handler
		staticAssetHandler := staticAssetServer()
		staticAssetHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
			staticAssetHandler.ServeHTTP(w, r)
		}
		http.HandleFunc("/assets/", headerMiddleware(staticAssetHandlerFunc))

		// Markdown File Handler
		markdownHandlerFunc := markdownHandleFunc(MarkdownHandlerOptions{
			DocRoot:       *dir,
			DocExtension:  *extension,
			DirIndex:      *index,
			MarkdownTheme: *markdownTheme,
			CodeTheme:     *codeTheme,
		})
		http.HandleFunc("/", basicAuthMiddleware(headerMiddleware(markdownHandlerFunc), *users))

		// Start HTTP server
		addr := fmt.Sprintf("%s:%d", *host, *port)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			fmt.Fprint(os.Stderr, "Error starting server (%s).", err)
			cli.Exit(1)
		}

	}

	app.Run(os.Args)
}
