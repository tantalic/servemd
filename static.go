package main

import "net/http"

func staticAssetServer() http.Handler {
	return http.StripPrefix("/assets", http.FileServer(assetFS()))
}
