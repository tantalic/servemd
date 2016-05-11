package main

import "net/http"

func staticAssetHandler() http.Handler {
	return http.StripPrefix("/assets", http.FileServer(assetFS()))
}
