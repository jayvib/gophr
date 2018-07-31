package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func HandleNewUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	RenderTemplate(w, r, "users/new", nil)
}
