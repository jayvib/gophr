package main

import (
	"net/http"
	)

func HandleImageNew(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "images/new", nil)
}

func HandleImageCreate(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("url") != "" {
		HandleImageCreateFromURL(w, r)
	} else {
		HandleImageCreateFromFile(w, r)
	}
}

func HandleImageCreateFromURL(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	image := NewImage(user)
	image.Description = r.FormValue("description")
	err := image.CreateFromURL(r.FormValue("url"))
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "images/new", map[string]interface{}{
				"Error": err,
				"ImageURL": r.FormValue("url"),
				"Image": image,
			})
		}
		panic(err)
	}
	http.Redirect(w, r, "/?flash=Image+Uploaded+Successfully", http.StatusFound)
}

func HandleImageCreateFromFile(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	image := NewImage(user)
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	if file == nil {
		RenderTemplate(w, r, "images/new", map[string]interface{}{
			"Image": image,
			"Error": errNoImage,
		})
		return
	}
	defer file.Close()
	err = image.CreateFromFile(file, header)
	if err != nil {
		RenderTemplate(w, r, "images/new", map[string]interface{}{
			"Image": image,
			"Error": err,
		})
		return
	}
	http.Redirect(w, r, "/?flash=Image+Uploaded+Successfully", http.StatusFound)
}