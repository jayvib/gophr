package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func HandleNewUserPage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "users/new", nil)
}

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
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

func HandleUserEdit(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	RenderTemplate(w, r, "users/edit", map[string]interface{}{
		"User": user,
	})
}

func HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	currentUser := RequestUser(r)
	email := r.FormValue("email")
	currentPassword := r.FormValue("currentPassword")
	newPassword := r.FormValue("newPassword")

	user, err := UpdateUser(currentUser, email, currentPassword, newPassword)
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "users/edit", map[string]interface{}{
				"Error": err,
				"User": user,
			})
			return
		}
		panic(err)
	}
	err = globalUserStore.Save(*currentUser)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/private/account?flash=User+updated", http.StatusFound)
}