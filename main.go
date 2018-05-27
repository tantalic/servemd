package main

import (
	"fmt"
	"net/http"
	"os"

	"log"

	"github.com/jawher/mow.cli"
)

const (
	Version      = "0.7.1"
	DefaultTheme = "clean"
)

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
		robotsTag = app.String(cli.StringOpt{
			Name:   "r x-robots-tag",
			Desc:   "Sets a X-Robots-Tag header",
			EnvVar: "X_ROBOTS_TAG",
			Value:  "",
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
			Value:  DefaultTheme,
			EnvVar: "MARKDOWN_THEME",
		})
		typekitKitID = app.String(cli.StringOpt{
			Name:   "t typekit-kit-id",
			Desc:   "ID of webfont kit to include from typekit",
			Value:  DefaultTheme,
			EnvVar: "TYPEKIT_KIT_ID",
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
		staticAssetHandlerFunc = headerMiddleware(staticAssetHandlerFunc)
		staticAssetHandlerFunc = basicAuthMiddleware(staticAssetHandlerFunc, *users)
		staticAssetHandlerFunc = robotsTagMiddleware(staticAssetHandlerFunc, *robotsTag)
		http.HandleFunc("/assets/", staticAssetHandlerFunc)

		// Setup the markdown theme (may be custom or bundled)
		themePath, themeHandler := theme(*markdownTheme)
		if themeHandler != nil {
			themeHandler = headerMiddleware(themeHandler)
			themeHandler = basicAuthMiddleware(themeHandler, *users)
			themeHandler = robotsTagMiddleware(themeHandler, *robotsTag)
			http.HandleFunc(themePath, themeHandler)
		}

		// Markdown File Handler
		markdownHandlerFunc := markdownHandleFunc(MarkdownHandlerOptions{
			DocRoot:       *dir,
			DocExtension:  *extension,
			DirIndex:      *index,
			MarkdownTheme: themePath,
			TypekitKitID:  *typekitKitID,
			CodeTheme:     *codeTheme,
		})
		markdownHandlerFunc = headerMiddleware(markdownHandlerFunc)
		markdownHandlerFunc = basicAuthMiddleware(markdownHandlerFunc, *users)
		markdownHandlerFunc = robotsTagMiddleware(markdownHandlerFunc, *robotsTag)
		http.HandleFunc("/", markdownHandlerFunc)

		// Start HTTP server
		addr := fmt.Sprintf("%s:%d", *host, *port)
		log.Printf("Starting server on %s", addr)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Printf("Error starting server: %s", err)
			cli.Exit(1)
		}

	}

	app.Run(os.Args)
}
