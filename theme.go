package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func theme(value string) (string, http.HandlerFunc) {
	// Check if the value is a path to a CSS file
	if filepath.Ext(value) == ".css" {
		cssPath, err := filepath.Abs(value)
		if err != nil {
			log.Printf("Error finding custom theme: %s: %s. Falling back to: %s.", value, err.Error(), DefaultTheme)
			return themePath(DefaultTheme), nil
		}

		if stat, err := os.Stat(cssPath); err == nil && !stat.IsDir() {
			// Add handler to serve content
			handler := func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, cssPath)
			}
			return themePath("custom"), handler
		}

		log.Printf("Error finding custom theme: %s. Falling back to: %s.", value, DefaultTheme)
		return themePath(DefaultTheme), nil
	}

	// If the value is not a path, check that it is a valid bundled theme
	localPath := themeAssetPath(value)
	_, err := Asset(localPath)
	if err != nil {
		log.Printf("Invalid theme selected: %s. Falling back to: %s.", value, DefaultTheme)
		return themePath(DefaultTheme), nil
	}

	return themePath(value), nil
}

func themePath(value string) string {
	return fmt.Sprintf("/assets/themes/%s.css", value)
}

func themeAssetPath(value string) string {
	return fmt.Sprintf("assets/themes/%s.css", value)
}
