package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Content
	DocRoot      string `envconfig:"DOCUMENT_ROOT" default:"."`
	DocExtension string `envconfig:"DOCUMENT_EXTENSION" default:".md"`
	DirIndex     string `envconfig:"DIRECTORY_INDEX" default:"index"`

	// Server
	Host string `envconfig:"HOST" default:""`
	Port int    `envconfig:"PORT" default:"3000"`

	// Theme
	MarkdownTheme string `envconfig:"MARKDOWN_THEME" default:""`
	CodeTheme     string `envconfig:"CODE_THEME" default:""`
}

type httpHandleFunc func(w http.ResponseWriter, r *http.Request)

func main() {
	// Load cnfiguation (environment variables)
	config, err := getConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration from environment (%s).\n", err)
		os.Exit(1)
	}

	// Static Asset Handler
	http.Handle("/assets/", staticAssetHandler(config))

	// Markdown File Handler
	http.HandleFunc("/", markdownHandleFunc(config))

	// Start HTTP server
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server (%s).\n", err)
		os.Exit(1)
	}
}

func getConfig() (Config, error) {
	var c Config
	err := envconfig.Process("MDSERVE", &c)
	return c, err
}
