package main

import (
	"net/http"
	)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index/home", nil)
}
