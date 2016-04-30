package main

import "net/http"

func staticAssetHandler(c Config) http.Handler {
	return http.StripPrefix("/assets", http.FileServer(assetFS()))
}
