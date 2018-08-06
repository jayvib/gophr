package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func HandleNewUserPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	RenderTemplate(w, r, "users/new", nil)
}

func HandleCreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := NewUser(
		r.FormValue("username"),
		r.FormValue("email"),
		r.FormValue("password"),
	)
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "users/new", map[string]interface{}{
				"Error": err,
				"User": user,
			})
		}
		panic(err)
	}
	http.Redirect(w, r, "/?flash=User+created", http.StatusFound)
}

func HandleLoginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	
}