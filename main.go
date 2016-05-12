package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
)

const (
	Version = "0.3.1"
)

type httpHandleFunc func(w http.ResponseWriter, r *http.Request)

func main() {

	app := cli.NewApp()

	// App Info
	app.Name = "servemd"
	app.Usage = "a simple http server for markdown content"
	app.UsageText = app.Name + " [options]"
	app.Version = Version
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Kevin Stock",
			Email: "kevinstock@tantalic.com",
		},
	}

	// CLI Flags
	app.Flags = []cli.Flag{
		// HTTP Server
		cli.StringFlag{
			Name:   "host",
			Value:  "0.0.0.0",
			Usage:  "the host/ip address to listen on for http",
			EnvVar: "HOST",
		},
		cli.IntFlag{
			Name:   "port",
			Value:  3000,
			Usage:  "the port to listen on for http",
			EnvVar: "PORT",
		},

		// Content
		cli.StringFlag{
			Name:   "dir",
			Value:  ".",
			Usage:  "the content directory to serve",
			EnvVar: "DOCUMENT_ROOT",
		},
		cli.StringFlag{
			Name:   "extension",
			Value:  ".md",
			Usage:  "the extension used for markdown files",
			EnvVar: "DOCUMENT_EXTENSION",
		},
		cli.StringFlag{
			Name:   "index",
			Value:  "index",
			Usage:  "the filename (without extension) to use for directory index",
			EnvVar: "DIRECTORY_INDEX",
		},

		// Theme
		cli.StringFlag{
			Name:   "markdown-theme",
			Value:  "clean",
			Usage:  "the theme to use for styling markdown html",
			EnvVar: "MARKDOWN_THEME",
		},
		cli.StringFlag{
			Name:   "code-theme",
			Usage:  "the highlight.js theme to use for syntax highlighting",
			EnvVar: "CODE_THEME",
		},
	}

	app.Action = start
	app.Run(os.Args)
}

func start(c *cli.Context) error {
	// Static Asset Handler
	http.Handle("/assets/", staticAssetHandler())

	// Markdown File Handler
	http.HandleFunc("/", markdownHandleFunc(MarkdownHandlerOptions{
		DocRoot:       c.String("dir"),
		DocExtension:  c.String("extension"),
		DirIndex:      c.String("index"),
		MarkdownTheme: c.String("markdown-theme"),
		CodeTheme:     c.String("code-theme"),
	}))

	// Start HTTP server
	addr := fmt.Sprintf("%s:%d", c.String("host"), c.Int("port"))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return fmt.Errorf("Error starting server (%s).", err)
	}

	return nil
}
