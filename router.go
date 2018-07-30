package main

import "github.com/julienschmidt/httprouter"

func NewRouter() *httprouter.Router {
	return httprouter.New()
}
