package main

import (
	"net/http"

	"github.com/gorilla/mux"
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
				"User":  user,
			})
		}
		panic(err)
	}
	http.Redirect(w, r, "/?flash=User+created", http.StatusFound)
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
				"User":  user,
			})
			return
		}
		panic(err)
	}
	err = globalUserStore.Save(*currentUser)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/account?flash=User+updated", http.StatusFound)
}

func HandleUserShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userID"]

	user, err := globalUserStore.Find(userId)
	if err != nil {
		panic(err)
	}
	if user == nil {
		http.NotFound(w, r)
		return
	}

	images, err := globalImageStore.FindAllByUser(user, 0)
	if err != nil {
		panic(err)
	}
	RenderTemplate(w, r, "users/show", map[string]interface{}{
		"Images": images,
		"User":   user,
	})
}
