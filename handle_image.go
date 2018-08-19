package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
				"Error":    err,
				"ImageURL": r.FormValue("url"),
				"Image":    image,
			})
		}
		panic(err)
	}
	http.Redirect(w, r, "/?flash=Image+Uploaded+Successfully", http.StatusFound)
}

// HandleImageCreateFromFile is a handler for saving the image file which the use uploaded.
func HandleImageCreateFromFile(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	image := NewImage(user)
	image.Description = r.FormValue("description")
	fmt.Printf("Image Description: %s\n", image.Description)
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

// HandleImageShow shows the image that uploaded by the user.
func HandleImageShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := vars["imageID"]
	image, err := globalImageStore.Find(imageID)
	if err != nil {
		panic(err)
	}
	if image == nil {
		http.NotFound(w, r)
		return
	}

	user, err := globalUserStore.Find(image.UserID)
	if err != nil {
		panic(err)
	}
	if user == nil {
		panic(fmt.Errorf("Could not find user %s", image.UserID))
	}
	RenderTemplate(w, r, "image/show", map[string]interface{}{
		"Image": image,
		"User":  user,
	})
}
